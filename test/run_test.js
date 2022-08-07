import "./polyfill_performance.js";
import { Go } from "./wasm_exec.js";

const go = new Go();

const wasm = await Deno.readFile(Deno.args[1]);
const result = await WebAssembly.instantiate(wasm)
go.run(result.instance);
