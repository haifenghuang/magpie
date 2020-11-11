## Run magpie programs in the browser using WebAssembly

In the magpie home directory, just run `run_wasm.sh`, then
open your browser, and type `http://localhost:9090`.


Below is the snapshot of the running magpie demo in browser:

<p>
    <img alt="magpie playground" src="https://github.com/haifenghuang/magpie/blob/master/wasm/magpie_playground.png?raw=true" width="450" height="450">
</p>

<br>
<p>
    <img alt="magpie playground" src="https://github.com/haifenghuang/magpie/blob/master/wasm/magpie_playground2.png?raw=true" width="450" height="450">
</p>

<br>
<p>
    <img alt="magpie playground" src="https://github.com/haifenghuang/magpie/blob/master/wasm/magpie_playground3.png?raw=true" width="450" height="450">
</p>

## Limitation

1. Can not use 'stdin', 'stdout' and 'stderr'
2. Can not use file object's method. e.g. fileObj.read(), fileObj.readLine()