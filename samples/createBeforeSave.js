function beforeSave(input) {
  input.concat = input.name + " " + input.score;
  return input;
}

module.exports = beforeSave;
