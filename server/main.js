import { Meteor } from 'meteor/meteor';

Meteor.startup(function () {
  // Load future from fibers
  var Future = Npm.require("fibers/future");
  // Load exec
  var exec = Npm.require("child_process").exec;

  // Server methods
  Meteor.methods({

    SequenceDownloader: function (retmax, term, taxon, minSeqLen, maxSeqLen, addlFlags) {

      // This method call won't return immediately, it will wait for the
      // asynchronous code to finish, so we call unblock to allow this client
      // to queue other method calls (see Meteor docs)

      this.unblock();

      var future=new Future();

      if( retmax == ' -retmax ' ) { // If -retmax is not set, set to default retmax=1
        retmax=" -retmax 1"
      }

      if( term == ' -term ' ) { // If term flag is not set don't include the -term flag in the call
        term=""
      }

      if( taxon == ' -taxon ') { // If taxon flag is not set don't include the -taxon flag in the call
        taxon=""
      }
	if( minSeqLen == ' -minSeqLen '){
		minSeqLen = ' -minSeqLen 0 '
	}
	if( maxSeqLen == ' -maxSeqLen '){
		maxSeqLen = ' -maxSeqLen 1000000000000000000000000000000000000000000000 '
	}



      // for local deployment
      var command='/Users/mathewsimenc/dropbox_csuf/Eernisse_White_ncbi_sequence-dler/web/sequencedownloader/private/ncbi-output' + retmax + taxon + term + minSeqLen + maxSeqLen + addlFlags; // Example term flag: -term title,mitochondrion,AND
      // for web deployment
      //var command='../ncbi-output-linux' + retmax + taxon + term + minSeqLen + maxSeqLen + addlFlags; // Example term flag: -term title,mitochondrion,AND

//      var fs = require('fs');
//      var files = fs.readdirSync('../');
//      future.return(files)
//      return future.wait();
      console.log(command);

      exec(command, {maxBuffer: 1024 * 50000}, function(error,stdout,stderr){

        if(error){

          console.log(error);

          throw new Meteor.Error(500,command+" failed");

        }

        future.return(stdout.toString());

      });

      return future.wait();

    }

  });
});
