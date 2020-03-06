const customLogic = require("./customLogic.js");
const fs = require("fs");

const input = JSON.parse(fs.readFileSync("./input.json"));
const output = customLogic(input);
fs.writeFileSync("./output/output.json", JSON.stringify(output));
