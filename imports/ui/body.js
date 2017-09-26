import { Template } from 'meteor/templating';

import './body.html'; 

import { saveAs } from 'file-saver';

import { TweenMax, TimelineMax } from 'gsap'; 


const parseSearchParams = function(numRecords=false) {

	event.preventDefault();
	const target = event.target;
	let taxon = ' -taxon ' + target.taxon.value.split(' ').join('+')
	let retmax = ' -retmax ' + target.retmax.value;
	let minSeqLen = ' -minSeqLen ' + target.minSeqLength.value;
	let maxSeqLen = ' -maxSeqLen ' + target.maxSeqLength.value;
	let term = '' // Will contain all -term flags to be passed to executable
	let addlFlags = ''

	let atp6 = target.ATP6; // Pre-made search term checkboxes
	let atp8 = target.ATP8;
	let coi = target.COI;
	let coii = target.COII;
	let coiii = target.COIII;
	let cyb = target.CYB;
	let nad1 = target.NAD1;
	let nad2 = target.NAD2;
	let nad3 = target.NAD3;
	let nad4 = target.NAD4;
	let nad4l = target.NAD4L;
	let nad5 = target.NAD5;
	let nad6 = target.NAD6;
	let rdna16S = target.rdna16S;
	let rdna12S = target.rdna12S;

	let outputFilename = target.taxon.value

	//alert(window.location.pathname);

	//let term_term = target.term.value.split(' ').join('+')
	let numRows = Session.get('currentNumRows');

	for (var i=1; i<=numRows; i++) { // Collect all terms
		let thisTermRow = 'term' + i;
		let thisFieldRow = 'ncbifield' + i;
		let thisLogicRow = 'andornot' + i;
		thisTerm = document.getElementById(thisTermRow).value.split(' ').join('+')
		if ( !!thisTerm ) {
			thisField = document.getElementById(thisFieldRow).value.split(' ').join('+')
			thisLogic = document.getElementById(thisLogicRow).value.split(' ').join('+')
			term += ' -term ' + thisField + ',' + thisTerm + ',' + thisLogic;
			outputFilename += '_' + thisTerm
		}
	}


	if ( $(atp6).prop('checked') ) { // Create flags for executable based on pre-made search term boxes checked or unchecked
		term += ' -term Gene+Name,ATP6,OR'
		outputFilename += '_ATP6'
	}

	if ( $(atp8).prop('checked') ) {
		term += ' -term Gene+Name,ATP8,OR'
		outputFilename += '_ATP8'
	}

	if ( $(coi).prop('checked') ) {
		term += ' -term Gene+Name,COI,OR -term Gene+Name,COXI,OR -term Gene+Name,COX1,OR -term Gene+Name,CO1,OR'
		outputFilename += '_COI'
	}

	if ( $(coii).prop('checked') ) {
		term += ' -term Gene+Name,COII,OR -term Gene+Name,COXII,OR -term Gene+Name,COX2,OR -term Gene+Name,CO2,OR'
		outputFilename += '_COII'
	}

	if ( $(coiii).prop('checked') ) {
		term += ' -term Gene+Name,COIII,OR -term Gene+Name,COXIII,OR -term Gene+Name,COX3,OR -term Gene+Name,CO3,OR'
		outputFilename += '_COIII'
	}

	if ( $(cyb).prop('checked') ) {
		term += ' -term Gene+Name,CYB,OR -term Gene+Name,CYTB,OR'
		outputFilename += '_CYB'
	}

	if ( $(nad1).prop('checked') ) {
		term += ' -term Gene+Name,ND1,OR -term Gene+Name,NAD1,OR'
		outputFilename += '_ND1'
	}

	if ( $(nad2).prop('checked') ) {
		term += ' -term Gene+Name,ND2,OR -term Gene+Name,NAD2,OR'
		outputFilename += '_ND2'
	}

	if ( $(nad3).prop('checked') ) {
		term += ' -term Gene+Name,ND3,OR -term Gene+Name,NAD3,OR'
		outputFilename += '_ND3'
	}

	if ( $(nad4).prop('checked') ) {
		term += ' -term Gene+Name,ND4,OR -term Gene+Name,NAD4,OR'
		outputFilename += '_ND4'
	}

	if ( $(nad4l).prop('checked') ) {
		term += ' -term Gene+Name,ND4L,OR -term Gene+Name,NAD4L,OR'
		outputFilename += '_ND4L'
	}

	if ( $(nad5).prop('checked') ) {
		term += ' -term Gene+Name,ND5,OR -term Gene+Name,NAD5,OR'
		outputFilename += '_ND5'
	}

	if ( $(nad6).prop('checked') ) {
		term += ' -term Gene+Name,ND6,OR -term Gene+Name,NAD6,OR'
		outputFilename += '_ND6'
	}

	if ( $(rdna16S).prop('checked') ) {
		term += ' -term All+Fields,16S+ribosomal+RNA,OR'
		outputFilename += '_16S'
	}

	if ( $(rdna12S).prop('checked') ) {
		term += ' -term All+Fields,12S+ribosomal+RNA,OR'
		outputFilename += '_12S'
	}

	if ( term == '' ) {
		term = ' -term '
	} else {
		if ( $(completeGenome).prop('checked') ) {
			addlFlags += ' -mito'
		}
		if ( $(nonGenome).prop('checked') ){
			addlFlags += ' -reg'
		}
	}

	//alert(term)


	//const outputFilename = target.taxon.value + '_' + term_term + '.csv'
	outputFilename = outputFilename.split('+').join('_')
	outputFilename += '.csv'

	return {outName: outputFilename, term: term, retmax: retmax, taxon: taxon, addlFlags: addlFlags, minSeqLen: minSeqLen, maxSeqLen: maxSeqLen}
	
}

	
Template.body.helpers({
	searchTermHelper: function() {
		return Session.get('searchTermRows');
	},

});

Template.body.events({

//	'submit #maxRecordsBtn': function() {
//
//		searchParams = parseSearchParams(numRecords=true)
//		outputfilename = searchParams["outName"]
//		term = searchParams["term"]
//		retmax = searchParams["retmax"]
//		taxon = searchParams["taxon"]
//
////		 Meteor.call('SequenceDownloader', retmax, term, taxon, function (err, response) {
////		//	var numRecordsFound = response
////		//	alert(numRecordsFound)
////			alert('hi')
////		 };
//	},


	'click .add-btn': function() {
		let numRows = Session.get('searchTermRows');
		let nextRow = Session.get('currentNumRows') + 1;
		numRows.push({row: nextRow});
		Session.set('currentNumRows', nextRow);
		Session.set('searchTermRows', numRows);
	},

	'click .delete-btn': function() {
		let numRows = Session.get('currentNumRows');
		let currentRow = event.target.value;
		let rowTerms = []
		let rowFields = []
		let rowLogics = []

		if (numRows == 1) {
			let thisTermRow = 'term1'
			let thisFieldRow = 'ncbifield1'
			let thisLogicRow = 'andornot1'
			document.getElementById(thisTermRow).value = ''
			document.getElementById(thisFieldRow).value = 'All+Fields'
			document.getElementById(thisLogicRow).value = 'AND'
			return

		}

		for (var i=1; i<=numRows; i++) { // Collect all values except the one deleted to repopulate rows

			if (i == currentRow ) {
				continue
			}

			let thisTermRow = 'term' + i;
			let thisFieldRow = 'ncbifield' + i;
			let thisLogicRow = 'andornot' + i;
			thisTerm = document.getElementById(thisTermRow).value
			thisField = document.getElementById(thisFieldRow).value
			thisLogic = document.getElementById(thisLogicRow).value
			rowTerms.push(thisTerm);
			rowFields.push(thisField);
			rowLogics.push(thisLogic);
		}

		newTermRows = []
		newNumRows = rowTerms.length

		for (var i=1; i<=newNumRows; i++) {
			newTermRows.push({row: i})
		}

		Session.set('searchTermRows', newTermRows);
		Session.set('currentNumRows', newNumRows);

		for (var i=1; i<=newNumRows; i++) { // Collect all values except the one deleted to repopulate rows

			let thisTermRow = 'term' + i;
			let thisFieldRow = 'ncbifield' + i;
			let thisLogicRow = 'andornot' + i;
			document.getElementById(thisTermRow).value = rowTerms[i-1]
			document.getElementById(thisFieldRow).value = rowFields[i-1]
			document.getElementById(thisLogicRow).value = rowLogics[i-1]

		}
	},

	'submit .search-params'(event) {

		var submitBtnId = document.activeElement.id

		searchParams = parseSearchParams()

		outputfilename = searchParams["outName"]
		term = searchParams["term"]
		retmax = searchParams["retmax"]
		taxon = searchParams["taxon"]
		addlFlags = searchParams["addlFlags"]
		minSeqLen = searchParams["minSeqLen"]
		maxSeqLen = searchParams["maxSeqLen"]

		if (submitBtnId == 'maxRecordsBtn'){
			addlFlags += ' -num'
		}

		 Meteor.call('SequenceDownloader', retmax, term, taxon, minSeqLen, maxSeqLen, addlFlags, function (err, response) {
			if (submitBtnId == 'maxRecordsBtn'){
				maxRecordInput = document.getElementById('retmax-entry')
				maxRecordInput.value = response
			} else {
				var file = new File([response], outputfilename, {type: "text/plain;charset=utf-8"});
				saveAs(file);
			}
      //		const filestream = streamsaver.createwritestream('filename.txt')
      //		const writer = filestream.getwriter()
      //		const encoder = new textencoder
      //		let output = encoder.encode(response)
      //		writer.write(response)
      //		writer.close()
		});

	}
});


Meteor.startup(function(){

	Session.set('searchTermRows', [{row: 1}]); // Set number of search term input lines
	Session.set('currentNumRows', 1);

//	$( function() { 
//		$("#seq-length-slider").slider({
//			range: true,
//			min: 0,
//			max: 100000,
//			values: [0,15000]
//		});
//	});


  
  //Create variables we will be referencing in our tweens.
  var white = 'rgb(255,255,255)';
  var seafoam = 'rgb(30, 129, 204, 1)';  
  //var seafoam = 'rgb(30,205,151)';  
  $buttonShapes = $('rect.btn-shape');
  $buttonColorShape = $('rect.btn-shape.btn-color');
  $buttonText = $('text.textNode');
  $buttonCheck = $('text.checkNode');
  
  //These are the button attributes which we will be tweening
  //This will be used with GSAP and the function below to tween
  var buttonProps = {
    buttonWidth : $buttonShapes.attr('width'),
    buttonX : $buttonShapes.attr('x'),
    buttonY : $buttonShapes.attr('y'),
    textScale : 1,
    textX : $buttonText.attr('x'),
    textY : $buttonText.attr('y')
  };
  
  //This is the update handler that lets us tween attributes
  function onUpdateHandler(){
    $buttonShapes.attr('width', buttonProps.buttonWidth);
    $buttonShapes.attr('x', buttonProps.buttonX);
    $buttonShapes.attr('y', buttonProps.buttonY);
    $buttonText.attr('transform', "scale(" + buttonProps.textScale + ")");
    $buttonText.attr('x', buttonProps.textX);
    $buttonText.attr('y', buttonProps.textY);
  }
  
  //Finally, create the timelines
  var hover_tl = new TimelineMax({
    tweens:[
      TweenMax.to( $buttonText, .15, { fill:white } ),
      TweenMax.to( $buttonShapes, .25, { fill: seafoam })
    ]
  });
  hover_tl.stop();
  
  var tl = new TimelineMax({onComplete:bind_mouseenter});
  //This is the initial transition, from [submit] to the circle
  tl.append( new TimelineMax({
    align:"start",
    tweens:[
      TweenMax.to( $buttonText, .15, { fillOpacity:0 } ),
      TweenMax.to( buttonProps, .25, { buttonX: (190-64)/2, buttonWidth:64, onUpdate:onUpdateHandler } ),
      TweenMax.to( $buttonShapes, .25, { fill: white })
    ], 
    onComplete:function(){ 
      $buttonColorShape.css({
        'strokeDasharray':202,
        'strokeDashoffset':202
      });
    }
  }) );
  
  //The loading dasharray offset animation… 
  tl.append(TweenMax.to($buttonColorShape, 1.2, {
    strokeDashoffset:0, 
    ease:Quad.easeIn,
    onComplete:function(){
      //Reset these values to their defaults.
      $buttonColorShape.css({
        'strokeDasharray':453,
        'strokeDashoffset':0
      });
    }
  }));

  //The Finish - transition to check
  tl.append(new TimelineMax({
    align:"start",
    tweens:[
      TweenMax.to($buttonShapes, .3, {fill:seafoam}),
      TweenMax.to( $buttonCheck, .15, { fillOpacity:1 } ),
      TweenMax.to( buttonProps, .25, { buttonX: 3, buttonWidth:190, onUpdate:onUpdateHandler } )
    ]
  }));
  
  //The Reset - back to the beginning
  //For demo only - probably you would want to remove this.
  tl.append(TweenMax.to($buttonCheck, .1, {delay:1,fillOpacity:0}));

  tl.append(new TimelineMax({
    align:"start",
    tweens:[
      TweenMax.to($buttonShapes, .3, {fill:white}),
      TweenMax.to($buttonText, .3, {fill:seafoam, fillOpacity:1})
    ],
    onComplete:function() {
      $('.colins-submit').removeClass('is-active');
    }
  }));
  tl.stop();
  
  //-- On Click, we launch into the cool transition
  $('.colins-submit').on('click', function(e) {
    //-- Add this class to indicate state
    $(e.currentTarget).addClass('is-active');
    tl.restart();
    $('.colins-submit').off('mouseenter');
    $('.colins-submit').off('mouseleave');
  });
  
  bind_mouseenter();
  
  function bind_mouseenter() {
    $('.colins-submit').on('mouseenter', function(e) {
      hover_tl.restart();
      $('.colins-submit').off('mouseenter');
      bind_mouseleave();
    });
  }
  function bind_mouseleave() {  
    $('.colins-submit').on('mouseleave', function(e) {
      hover_tl.reverse();
      $('.colins-submit').off('mouseleave');
      bind_mouseenter();      
    });
  }
  
});
