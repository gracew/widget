function beforeSave(input) {
  input.concat = input.name + " " + input.score;
}

module.exports = beforeSave;
