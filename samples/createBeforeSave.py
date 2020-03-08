def before_save(input):
  input["concat"] = input["name"] + " " + str(input["score"])
  return input
