def before_create(input):
  input["concat"] = input["name"] + " " + str(input["score"])
  return input
