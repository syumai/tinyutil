import "./polyfill_performance.js";
import { Go } from "./wasm_exec.js";

const go = new Go();

const wasm = await Deno.readFile(Deno.args[0]);
const result = await WebAssembly.instantiate(wasm, go.importObject);
go.run(result.instance);
