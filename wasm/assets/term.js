(function() {
	// Write output to HTML element.
	window.set_output = (output) => {
		// Write stdout to terminal.
		let outputBuf = '';
		const decoder = new TextDecoder("utf-8");
		global.fs.writeSync = (fd, buf) => {
			outputBuf += decoder.decode(buf);
			const nl = outputBuf.lastIndexOf("\n");
			if (nl != -1) {
				output.value += outputBuf.substr(0, nl + 1);
				//window.scrollTo(0, document.body.scrollHeight);
				outputBuf = outputBuf.substr(nl + 1);
			}
			return buf.length;
		};
	};

}());