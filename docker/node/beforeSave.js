const beforeSave = require("./customLogic.js");
const fs = require("fs");

const input = JSON.parse(fs.readFileSync("./input.json"));
beforeSave(input);
