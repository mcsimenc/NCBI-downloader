package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"net/http"
	"flag"
	"strings"
	"regexp"
)

func stripTag(input string) (output string) {

	// Returns value from XML line.
	splitString := strings.Split(input, ">")
	splitString2 := strings.Split(splitString[1], "<")
	output = splitString2[0]

	return

}

func findTag(lines []string, tag string, offset int) (output string) {

	// Finds first occurrence of <tag> in lines (XML) and returns value
	// from line that is offset number of tags beyond <tag>.

	output = "NA" // Default output (i.e. if tag is not found, output is NA)
	offsetCount := -1 // Counts the number of tags after the tag specified in args during loop
	tagRE, _ := regexp.Compile("<.*?>") // Regular expression matches any opening or closing tag

	for _, line := range lines {

		if offsetCount > -1 { // This will be true once the tag has been found, and will increment for each following tag. There are sometimes empty lines between tags which need to be ignored.
			if tagRE.MatchString(line) {
				offsetCount += 1
			}
		}

		if strings.Contains(line, tag) {
			offsetCount += 1
		}

		if offsetCount == offset {

			output = strings.Replace(stripTag(line), ",", "_", -1)
			break
		}
	}

	return
}

func findTags(lines []string, tag string) (output string) {

	// Finds first occurrence of <tag> in lines (XML format) and returns
	// semicolon-concatenated values for all sub-tags. Written with 
	// <GBReference_authors> in mind for which there are multiple <GBAuthor> subtags

	end_tag := strings.Replace(tag, "<", "</", 1)
	collect := false
	output = "NA"

	for _, line := range lines {

		if collect == true {

			val := stripTag(line)

			if output != "" {

				output = fmt.Sprintf("%s;%s", output, val)

			} else {

				output = val
			}
		}

		if strings.Contains(line, end_tag) {

			break

		} else if strings.Contains(line, tag) {

			collect = true
		}
	}

	return
}


func esearchString(retmax int, taxon string, terms map[int][]string) (concat_string string) {

	// Parses parameters (retmax and taxon are required) into an eutils URL.
	// Optional map terms contains key:[value, logic], e.g. title:[mitochondrion, AND]
	// which becomes +AND+mitochondrion[title] in the eutils URL
	// EXAMPLE: https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=nuccore&retmax=10&term=mollusca[organism]+AND+complete[title]+AND+genome[title]+AND+mitochondrion[title]

	concat_string = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=nuccore&retmax="
	concat_string += fmt.Sprint(retmax, "&term=AND+", taxon, "[organism]+(")

	for key := range terms {

		concat_string += fmt.Sprint("+", terms[key][2], "+", terms[key][1], "[", terms[key][0], "]")
	}

	concat_string += ")"

	return
}

type termsFlags []string // This will be implemented as type flag.Value using the following methods:

func (i *termsFlags) String() string { // String() method for type termsFlags (when implemented as type flag.Value in flag.Var call in main() )

	return ""
}

func (i *termsFlags) Set(value string) error { // Here the -term flags are actually parsed. These methods are run automatically when flag.Var is called during main() because type flag.Value's methods are run upon parsing by default

	*i = append(*i, strings.TrimSpace(value))

	return nil

}

func cleanBinom(binom string) (output string) {

	// Removes special characters from binom according to the following rules:
	// 1. "sp. " --> "sp"
	// 2. "cf. " --> "cf"
	// 3. "-"    --> "_"
	// 4. "."    --> "_"
	// 5. ":"    --> "_"
	// 6. " "    --> "_"

	output = strings.Replace(binom, "sp. ", "sp", -1)
	output = strings.Replace(binom, "cf. ", "cf", -1)
	output = strings.Replace(binom, "-", "_", -1)
	output = strings.Replace(binom, ".", "_", -1)
	output = strings.Replace(binom, ":", "_", -1)
	output = strings.Replace(binom, " ", "_", -1)

	return
}

func parseLocality(GB_country string) (country, locality string) {
	// Parses GenBank country tag value into country and locality. Example:
	// USA: Virginia --> USA, Virginia as country, locality
	splitStr := strings.SplitN(GB_country, ":", 2)

	if len(splitStr) == 0 {
		country = "NA"
		locality = "NA"

	} else if len(splitStr) == 1 {
		country = splitStr[0]
		locality = "NA"

	} else {
		country = splitStr[0]
		locality = strings.TrimLeft(splitStr[1], " :.-_")
	}

	return
}


func parseNonMitoGenomeXML(xmlLines []string, geneNames map[string]bool, minSeqLen int, maxSeqLen int) {

		GB_locus_id := findTag(xmlLines, "GBSeq_locus", 0)
		GB_seq_length := findTag(xmlLines, "GBSeq_length", 0)
		GB_strandedness := findTag(xmlLines, "GBSeq_strandedness", 0)
		GB_moltype := findTag(xmlLines, "GBSeq_moltype", 0)
		GB_toplogy := findTag(xmlLines, "GBSeq_topology", 0)
		GB_division := findTag(xmlLines, "GBSeq_division", 0)
		GB_update_date := findTag(xmlLines, "GBSeq_update-date", 0)
		GB_create_date := findTag(xmlLines, "GBSeq_create-date", 0)
		GB_definition := findTag(xmlLines, "GBSeq_definition", 0)
		GB_primary_accession := findTag(xmlLines, "GBSeq_primary-accession", 0)
		GB_accession_version := findTag(xmlLines, "GBSeq_accession-version", 0)
		GB_source := findTag(xmlLines, "GBSeq_source", 0)
		GB_organism := cleanBinom(findTag(xmlLines, "GBSeq_organism", 0)) // cleanBinom()
		GB_taxonomy := findTag(xmlLines, "GBSeq_taxonomy", 0)
		GB_taxon_id := strings.SplitAfter(findTag(xmlLines, "<GBQualifier_value>taxon:", 0), "taxon:")[1]
		GB_gene := findTag(xmlLines, "<GBQualifier_name>gene", 1)
		GB_product := findTag(xmlLines, "<GBQualifier_name>product", 1)
		GB_codon_start := findTag(xmlLines, "<GBQualifier_name>codon_start", 1)
		GB_organelle := findTag(xmlLines, "<GBQualifier_name>organelle", 1)
		GB_pub_title := findTag(xmlLines, "<GBReference_title>", 0)
		GB_pub_jrn := findTag(xmlLines, "<GBReference_journal>", 0)
		GB_pub_authors := strings.TrimRight(strings.Replace(findTags(xmlLines, "<GBReference_authors>"), ",", " ", -1), ";") // multiple tags, used findTags()
		GB_voucher := findTag(xmlLines, "<GBQualifier_name>specimen_voucher", 1)
		GB_country := findTag(xmlLines, "<GBQualifier_name>country", 1)
		country, locality := parseLocality(GB_country)
		GB_lat_long := findTag(xmlLines, "<GBQualifier_name>lat_long", 1)
		GB_note := findTag(xmlLines, "<GBQualifier_name>note", 1)
		GB_haplotype := findTag(xmlLines, "<GBQualifier_name>haplotype", 1)
		GB_bio_material := findTag(xmlLines, "<GBQualifier_name>bio_material", 1)
		GB_isolation_source := findTag(xmlLines, "<GBQualifier_name>isolation_source", 1)
		GB_pop_variant := findTag(xmlLines, "<GBQualifier_name>pop_variant", 1)
		GB_isolate := findTag(xmlLines, "<GBQualifier_name>isolate", 1)
		GB_comment := findTag(xmlLines, "<GBSeq_comment", 0)
		GB_prot_sequence := findTag(xmlLines, "<GBQualifier_name>translation", 1)
		GB_cds_sequence := findTag(xmlLines, "<GBQualifier_name>transcription", 1)
		GB_nuc_sequence := findTag(xmlLines, "GBSeq_sequence", 0)

		if len(GB_nuc_sequence) < minSeqLen || len(GB_nuc_sequence) > maxSeqLen {
			GB_nuc_sequence = "NA"
		}

		fastaHeader := ">" + GB_organism + "_" + GB_primary_accession
		fastaHeaderAndNucSeq := fastaHeader + "$" + GB_nuc_sequence
		fastaHeaderAndCDSSeq := fastaHeader + "$" + GB_cds_sequence
		fastaHeaderAndProtSeq := fastaHeader + "$" + GB_prot_sequence

		if len(geneNames) > 0 { // If "Gene Name" field was specified with search terms filter results that don't have that name as GB_gene
			if !geneNames[strings.ToLower(GB_gene)] { // Skip this record if GB_gene is not one of the terms specified with Gene+Name field and OR/AND logic

				return
			}
		}

		// Print CSV lines to stdout
		outputFields := []string {GB_locus_id,
					GB_seq_length,
					GB_strandedness,
					GB_moltype,
					GB_toplogy,
					GB_division,
					GB_update_date,
					GB_create_date,
					GB_definition,
					GB_primary_accession,
					GB_accession_version,
					GB_source,
					GB_organism,
					GB_taxonomy,
					GB_nuc_sequence,
					GB_cds_sequence,
					GB_prot_sequence,
					fastaHeader,
					fastaHeaderAndNucSeq,
					fastaHeaderAndCDSSeq,
					fastaHeaderAndProtSeq,
					GB_taxon_id,
					GB_gene,
					GB_product,
					GB_codon_start,
					GB_organelle,
					GB_pub_title,
					GB_pub_authors,
					GB_pub_jrn,
					GB_voucher,
					country,
					locality,
					GB_lat_long,
					GB_note,
					GB_haplotype,
					GB_bio_material,
					GB_isolation_source,
					GB_pop_variant,
					GB_isolate,
					GB_comment}

		fmt.Println(strings.Join(outputFields,","))
}

func parseMitoGenomeXML(xmlLines []string, geneNames map[string]bool) {

	GB_moltype := findTag(xmlLines, "GBSeq_moltype", 0) // DNA, refers to the whole genome
	GB_seq_length := findTag(xmlLines, "GBSeq_length", 0) // Refers to the whole genome
	GB_nuc_sequence := "NA" // findTag(xmlLines, "GBSeq_sequence", 0) // Not present in mito genomes usually
	GB_locus_id := findTag(xmlLines, "GBSeq_locus", 0)
	GB_strandedness := findTag(xmlLines, "GBSeq_strandedness", 0)
	GB_toplogy := findTag(xmlLines, "GBSeq_topology", 0)
	GB_division := findTag(xmlLines, "GBSeq_division", 0)
	GB_update_date := findTag(xmlLines, "GBSeq_update-date", 0)
	GB_create_date := findTag(xmlLines, "GBSeq_create-date", 0)
	GB_definition := findTag(xmlLines, "GBSeq_definition", 0)
	GB_primary_accession := findTag(xmlLines, "GBSeq_primary-accession", 0)
	GB_accession_version := findTag(xmlLines, "GBSeq_accession-version", 0)
	GB_source := findTag(xmlLines, "GBSeq_source", 0)
	GB_organism := cleanBinom(findTag(xmlLines, "GBSeq_organism", 0)) // cleanBinom()
	GB_taxonomy := findTag(xmlLines, "GBSeq_taxonomy", 0)
	GB_taxon_id := strings.SplitAfter(findTag(xmlLines, "<GBQualifier_value>taxon:", 0), "taxon:")[1]
	GB_organelle := findTag(xmlLines, "<GBQualifier_name>organelle", 1)
	GB_pub_title := findTag(xmlLines, "<GBReference_title>", 0)
	GB_pub_jrn := findTag(xmlLines, "<GBReference_journal>", 0)
	GB_pub_authors := strings.TrimRight(strings.Replace(findTags(xmlLines, "<GBReference_authors>"), ",", " ", -1), ";") // multiple tags, used findTags()
	GB_voucher := findTag(xmlLines, "<GBQualifier_name>specimen_voucher", 1)
	GB_country := findTag(xmlLines, "<GBQualifier_name>country", 1)
	country, locality := parseLocality(GB_country)
	GB_lat_long := findTag(xmlLines, "<GBQualifier_name>lat_long", 1)
	GB_note := findTag(xmlLines, "<GBQualifier_name>note", 1)
	GB_haplotype := findTag(xmlLines, "<GBQualifier_name>haplotype", 1)
	GB_bio_material := findTag(xmlLines, "<GBQualifier_name>bio_material", 1)
	GB_isolation_source := findTag(xmlLines, "<GBQualifier_name>isolation_source", 1)
	GB_pop_variant := findTag(xmlLines, "<GBQualifier_name>pop_variant", 1)
	GB_isolate := findTag(xmlLines, "<GBQualifier_name>isolate", 1)
	GB_comment := findTag(xmlLines, "<GBSeq_comment", 0)

	fastaHeader := ">" + GB_organism + "_" + GB_primary_accession
	fastaHeaderAndNucSeq := fastaHeader + "$" + GB_nuc_sequence
	fastaHeaderAndCDSSeq := "NA"
	fastaHeaderAndProtSeq := "NA"

	// One per gene
	GB_gene := "NA"
	GB_product := "NA"
	GB_codon_start := "NA"
	GB_prot_sequence := "NA"
	GB_cds_sequence := "NA"

	foundGeneTag := false
	foundGene := false
	foundCodonStartTag := false
	foundProductTag := false
	foundTranslationTag := false

	for _, line := range xmlLines {

		if strings.Contains(line, "<GBQualifier_value>") && foundGeneTag {
			foundGeneTag = false
			gene := strings.Replace(stripTag(line), ",", "_", -1)

			if geneNames[strings.ToLower(gene)] { // Check if this is one of the genes searched for
				GB_gene = gene
				foundGene = true
			}
		}

		if strings.Contains(line, "<GBQualifier_name>gene") { // prints when encounters next gene tag or end of feature table section after which there should be no more genes

			if strings.Contains(line, "<GBQualifier_name>gene") {
				foundGeneTag = true
			}
		}

		if foundGene && strings.Contains(line, "</GBFeature>") {
			foundGene = false
			foundCodonStartTag = false
			foundProductTag = false
			foundTranslationTag = false

			if GB_prot_sequence == "NA" { // GBFeature tag may have gene tag but no sequence information. Don't write if that's the case
				GB_gene = "NA"
				GB_product = "NA"
				GB_codon_start = "NA"
				continue
			}

			// Print CSV lines to stdout
			outputFields := []string {GB_locus_id,
						GB_seq_length,
						GB_strandedness,
						GB_moltype,
						GB_toplogy,
						GB_division,
						GB_update_date,
						GB_create_date,
						GB_definition,
						GB_primary_accession,
						GB_accession_version,
						GB_source,
						GB_organism,
						GB_taxonomy,
						GB_nuc_sequence,
						GB_cds_sequence,
						GB_prot_sequence,
						fastaHeader,
						fastaHeaderAndNucSeq,
						fastaHeaderAndCDSSeq,
						fastaHeaderAndProtSeq,
						GB_taxon_id,
						GB_gene,
						GB_product,
						GB_codon_start,
						GB_organelle,
						GB_pub_title,
						GB_pub_authors,
						GB_pub_jrn,
						GB_voucher,
						country,
						locality,
						GB_lat_long,
						GB_note,
						GB_haplotype,
						GB_bio_material,
						GB_isolation_source,
						GB_pop_variant,
						GB_isolate,
						GB_comment}

			fmt.Println(strings.Join(outputFields,","))

			GB_gene = "NA"
			GB_product = "NA"
			GB_codon_start = "NA"
			GB_prot_sequence = "NA"
		}

		if foundGene { // Runs if previous gene tag is a gene searched for
			if foundCodonStartTag {
				GB_codon_start = strings.Replace(stripTag(line), ",", "_", -1)
				foundCodonStartTag = false
			}

			if strings.Contains(line, "<GBQualifier_name>codon_start") {
				foundCodonStartTag = true
			}

			if foundProductTag {
				GB_product = strings.Replace(stripTag(line), ",", "_", -1)
				foundProductTag = false
			}

			if strings.Contains(line, "<GBQualifier_name>product") {
				foundProductTag = true
			}

			if foundTranslationTag {
				GB_prot_sequence = strings.Replace(stripTag(line), ",", "_", -1)
				fastaHeaderAndProtSeq = fastaHeader + "$" + GB_prot_sequence
				foundTranslationTag = false
			}

			if strings.Contains(line, "<GBQualifier_name>translation") {
				foundTranslationTag = true
			}
		}
	}
}


func EfetchPOSTrequest(gb_ids string) (xmlLines []string) {

	concat_request := fmt.Sprint("https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=nuccore&id=", gb_ids, "&rettype=gb&retmode=xml")
	hc := http.Client{}
	form := url.Values{}
	req, _ := http.NewRequest("POST", concat_request, strings.NewReader(form.Encode()))
	gb_response, _ := hc.Do(req)
	gb_data, _ := ioutil.ReadAll(gb_response.Body)
	xmlString := string(gb_data)
	xmlLines = strings.Split(xmlString, "\n")

	return
}

func processGBSeqXMLrecords(xmlLines []string, geneNames map[string]bool, minSeqLen int, maxSeqLen int) {

	var xmlOneRecord []string

	EORindex := 0 //EOR = End of record

	for i, line := range xmlLines {

		if strings.Contains(line, "</GBSeq>") {
			xmlOneRecord = xmlLines[EORindex:i]
			EORindex = i
			// Loop through XML get definition. If definition has 'mitochondri' 'genome' and 'partial' or 'complete' set mitoGenome flag
			// If mitoGenomeFlag, run parseMitoGenomeXML() else run parseNonMitoGenomeXML()
			definition := findTag(xmlOneRecord, "GBSeq_definition", 0)

			if (strings.Contains(definition, "mitochondri") && strings.Contains(definition, "genome")) && (strings.Contains(definition, "partial") || strings.Contains(definition, "complete")) {
				parseMitoGenomeXML(xmlOneRecord, geneNames)

			} else {
				parseNonMitoGenomeXML(xmlOneRecord, geneNames, minSeqLen, maxSeqLen)
			}
		}
	}
}


func main() {

	taxonPtr := flag.String("taxon", "", "a string") // Will be used with [Organism] flag in eutils URL
	retmaxPtr := flag.Int("retmax", 1, "an int") // The maximum number of records to return from entrez search (the first n (retmax) encountered in search result XML will be returned)
	minSeqLen := flag.Int("minSeqLen", 0, "an int") // The maximum number of records to return from entrez search (the first n (retmax) encountered in search result XML will be returned)
	maxSeqLen := flag.Int("maxSeqLen", 15000, "an int") // The maximum number of records to return from entrez search (the first n (retmax) encountered in search result XML will be returned)
	mitoSearch := flag.Bool("mito", false, "bool") // Include complete or partial mito genomes in the search
	regSearch := flag.Bool("reg", false, "bool") // Exclude complete or partial mito genomes in the search
	//fastaOut := flag.Bool("fastaOut", false, "bool") // Exclude complete or partial mito genomes in the search
	returnNumRecords := flag.Bool("num", false, "bool") // If true simply do the search and return how many records were found
	// If both mitoSearch and regSearch are specified do not restrict search
	var terms termsFlags // To collect terms (multiple -term flags may be used)
	flag.Var(&terms, "term", "comma-sep string: label,searchTerm,logic   e.g. title,mitochondrion,AND   multiple -term may be specified") // becomes +AND+mitochondrion[title] in the eutils esearch URL
	flag.Parse()
	termMap := make(map[int][]string) // Turn comma-sep strings passed with -term into map: label,term,logic -> { int:[label, term, logic] }
	geneNames := make(map[string]bool) // Keep a collection of gene names to check output against

	if len(terms) > 0 {

		for i,term := range terms {
			splitTerms := strings.Split(term, ",")
			label := splitTerms[0]
			term := splitTerms[1]
			logic := splitTerms[2]
			val := []string{ label, term, logic }
			termMap[i] = val

			if label == "Gene+Name" && logic != "NOT" { // Store search term if used with Gene Name field to filter out non-matching Esearch records

				geneNames[strings.ToLower(term)] = true
			}
		}
	}

	concat_string := esearchString(*retmaxPtr, *taxonPtr, termMap) // Assemble eutils esearch URL from command line params

	if *mitoSearch && !*regSearch {
		concat_string += "+AND+(complete+genome[All+Fields]+OR+partial+genome[All+Fields])"


	} else if *regSearch && !*mitoSearch {
		concat_string += "+NOT+complete+genome[All+Fields]+NOT+partial+genome[All+Fields])"
	}

	id_response, _ := http.Get(concat_string)
	htmlData, _ := ioutil.ReadAll(id_response.Body)
	htmlString := string(htmlData)
	splitString := strings.Split(htmlString, "\n") // Convert XML string into slice

	if *returnNumRecords { // Return number of records found, runs if -num specified
		countRE, _ := regexp.Compile("<Count>([0-9]+)</Count>")
		count := "NA"

		for _, line := range splitString {

			if strings.Contains(line, "<Count>") {
				count = countRE.FindStringSubmatch(line)[1]
				fmt.Println(count)

				return
			}
		}
	}

	fmt.Println(	`locus_id,seq_length,strandedness,moltype,toplogy,division,update_date,create_date,definition,primary_accession,accession_version,source,organism,taxonomy,nuc_sequence,cds_sequence,prot_sequence,fasta_header,fasta_nt,fasta_cds,fasta_prot,taxon_id,gene,product,codon_start,organelle,pub_title,pub_authors,pub_jrn,voucher,country,locality,lat_long,note,haplotype,bio_material,isolation_source,pop_variant,isolate,comment`) // write headers for fields

	gb_ids := ""
	num_ids := 0

	for _, line := range splitString { // Get up to 500 records in one XML at once
					   // Efetch with 600 worked but 1000 returned nothing in my testing
		if strings.Contains(line, "<Id>") {
			gb_id := stripTag(line)
			gb_ids = fmt.Sprintf("%s,%s", gb_ids, gb_id)
			num_ids ++
		}

		if num_ids == 500 {
			xmlLines := EfetchPOSTrequest(gb_ids)
			processGBSeqXMLrecords(xmlLines, geneNames, *minSeqLen, *maxSeqLen)
			gb_ids = ""
			num_ids = 0
		}
	}

	if gb_ids != "" {
			xmlLines := EfetchPOSTrequest(gb_ids)
			processGBSeqXMLrecords(xmlLines, geneNames, *minSeqLen, *maxSeqLen)
	}
}
