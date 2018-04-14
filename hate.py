from hatesonar import Sonar
import fire
import json

class Hate(object):
  """A simple calculator class."""

  def check(self, txt):
    sonar = Sonar()
    res = sonar.ping(text=txt)
    return json.dumps(res)

if __name__ == '__main__':
  fire.Fire(Hate)