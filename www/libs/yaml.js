"use strict";

/* 
	Taken from https://github.com/coolaj86/json2yaml.git
	
	Removed nodejs components to allow it to work in browser
*/
var YAML = {};

YAML.stringify = function(data) {

	var indentLevel = '';

	var typeOf = function(val) {
		var tStr = Object.prototype.toString.call(val)
			.split(/\s/)[1].replace(/\]$/, '').toLowerCase();
		//console.log(tStr);
		return tStr;
	}

	var handlers = {
		"undefined": function () {
			// objects will not have `undefined` converted to `null` as this may have unintended consequences
			// For arrays, however, this behavior seems appropriate
			return 'null';
		}, 
		"null": function () {
			return 'null';
		}, 
		"number": function (x) {
			return x;
		},
		"boolean": function (x) {
			return x ? 'true' : 'false';
		}, 
		"string": function (x) {
			// to avoid the string "true" being confused with the
			// the literal `true`, we always wrap strings in quotes
			return JSON.stringify(x);
		}, 
		"array": function (x) {
			var output = '';

			if (0 === x.length) {
				output += '[]';
				return output;
			}

			indentLevel = indentLevel.replace(/$/, '  ');
			
			x.forEach(function (y) {
				// TODO how should `undefined` be handled?
				var handler = handlers[typeOf(y)];
				if (!handler) {
				  throw new Error('what the crap: ' + typeOf(y));
				}
				output += '\n' + indentLevel + '- ' + handler(y);
			});

			indentLevel = indentLevel.replace(/  /, '');

			return output;
		}, 
		"object": function (x) {
			var output = '';

			if (0 === Object.keys(x).length) {
				output += '{}';
				return output;
			}

			indentLevel = indentLevel.replace(/$/, '  ');
		  	Object.keys(x).forEach(function (k) {
				var val = x[k], 
					handler = handlers[typeOf(val)];

				if ('undefined' === typeof val) {
					// the user should do delete obj.key and not obj.key = undefined but we'll error on the side of caution
					return;
			  	}

			  	if (!handler) {
					throw new Error('what the crap: ' + typeOf(val));
			  	}

			  	output += '\n' + indentLevel + k + ': ' + handler(val);
		  	});
		  	indentLevel = indentLevel.replace(/  /, '');

		  	return output;
		}, 
		"function": function () {
			// TODO this should throw or otherwise be ignored
			return '[object Function]';
		}
	}; // end handlers.

	return '---' + handlers[typeOf(data)](data) + '\n';
}
