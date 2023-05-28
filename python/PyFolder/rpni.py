import aalpy
from aalpy.learning_algs import run_RPNI

def RPNIlearn(s):

    data = []

    for sample in s["positive"]:
        new_str = sample.replace("0", "a")
        new_str = new_str.replace("1", "b")

        data.append((tuple(new_str), True))

    for sample in s["negative"]:
        new_str = sample.replace("0", "a")
        new_str = new_str.replace("1", "b")

        data.append((tuple(new_str), False))

    model = aalpy.learning_algs.run_RPNI(data, automaton_type='dfa')