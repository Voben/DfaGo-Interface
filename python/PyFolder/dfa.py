import json 

# expects a dfa object to export it as a JSON string
# ready to be used as an object in DfaGo
def ToDfaGoJSON(dfa):
    states_lst = dfa.states.copy()
    if len(dfa.states) > len(dfa.tf):
        states_lst.append("s")

    ret = {"Alphabet": [], "StartingState":0, "States":[], "dateCreated":"None",
        "depth":0, "dirty":False, "docType":"DfaGo/DFA", "version":1}
        
    for a in dfa.alphabet:
        ret["Alphabet"].append(int(a))

    ret["StartingState"] = dfa.states.index(dfa.start)

    for idx_s in range(0, len(states_lst)):
        state_dict = {}

        if states_lst[idx_s] in dfa.accs:
            state_dict["Label"] = 0
        else:
            state_dict["Label"] = 2

        next_lst = []

        if states_lst[idx_s] in dfa.tf:
            for v in dfa.tf[states_lst[idx_s]].values():
                next_lst.append(states_lst.index(v))
        else:
            next_lst = [len(states_lst) - 1, len(states_lst) - 1]
        
        state_dict["Next"] = next_lst
        
        state_dict["depth"] = 0
        state_dict["order"] = 0
        
        ret["States"].append(state_dict)

    return json.dumps(ret)

# obtain a DFA object in Python from a DfaGo JSON string
def DfaGo_String(jsonDFA_bytes):
    jsonDFA = json.loads(jsonDFA_bytes.decode())

    ret_DFA = DFA()
    
    for i in jsonDFA["Alphabet"]:
            ret_DFA.alphabet.append(str(i))

    ret_DFA.start = str(jsonDFA["StartingState"])

    for idx1 in range(0, len(jsonDFA["States"])):
        ret_DFA.states.append(str(idx1))

        if jsonDFA["States"][idx1]["Label"] == 0:
            ret_DFA.accs.append(str(idx1))

        for idx2 in range(0, len(jsonDFA["States"][idx1]["Next"])):
            if str(idx1) not in ret_DFA.tf:
                ret_DFA.tf[str(idx1)] = {}
            ret_DFA.tf[str(idx1)][str(idx2)] = str(jsonDFA["States"][idx1]["Next"][idx2])
        
    ret_DFA.states = ret_DFA.states[1:]  

    return ret_DFA

# will extract a DFA from a JSON file containing a DfaGo DFA
def from_DfaGo_file(filename):
    with open(filename) as f:
        d = json.load(f)

    ret_DFA = DfaGo_String(d)

    return ret_DFA  

# DFA class
class DFA:
    # init class that does not require any parameters
    def __init__(self, states = [], alphabet = [], start = "", tf = {}, accs = []):
        self.add_states(states)
        self.add_alphabet(alphabet)
        self.add_start(start)
        self.add_tf(tf)
        self.add_accepting(accs)

    # the following functions all populate the different attributes of the class
    # they are used in the init func
    def add_states(self, states):
        self.states = states

    def add_alphabet(self, alphabet):
        self.alphabet = alphabet

    def add_start(self, start):
        self.start = start

        if self.start not in self.states:
            self.states.append(self.start)

    def add_tf(self, tf):
        self.tf = tf

    def add_accepting(self, accs):
        self.accs = accs

        for acc in self.accs:
            if acc not in self.states:
                self.states.append(acc)

    # this function expects a string to be processed by the DFA
    # returns a boolean depending on if the word is part of the
    # DFA's language or not
    def label(self, string):
        curNode = self.start

        for char in string:
            
            if curNode not in self.tf:
                return False
            if char not in self.tf[curNode]:
                return False

            curNode = self.tf[curNode][char]

        if curNode in self.accs:
            return True
        else:
            return False

    # export a Dfa from Python to a dot file
    def export_dot(self, string):
        state_names = {}
        for state in self.states:
            state_names[state] = str(len(state_names) + 1)

        with open(string + ".dot", "w") as fp:
            fp.write("digraph G{\n")
            fp.write("0 [label=\"\", shape=point];\n")
            fp.write("0 -> 1;\n")

            for state in self.states:

                if state in self.accs:
                    nodeShape = "doublecircle"
                else:
                    nodeShape = "circle"
    
                fp.write(state_names[state] + " [label=\""+ state_names[state] + "\", shape=" + nodeShape + "];\n")

                if state in self.tf:
                    fp.write(state_names[state] + " -> " + state_names[self.tf[state][self.alphabet[0]]] + " [label=" + self.alphabet[0] +"];\n")
                    fp.write(state_names[state] + " -> " + state_names[self.tf[state][self.alphabet[1]]] + " [label=" + self.alphabet[1] +"];\n")
                else:
                    fp.write(state_names[state] + " -> " + state_names[state] + " [label=" + self.alphabet[0] +"];\n")
                    fp.write(state_names[state] + " -> " + state_names[state] + " [label=" + self.alphabet[1] +"];\n")

            fp.write("}")

    # export a JSON tring 
    def export_JSON(self):
        DFA_dict = {"states": str(self.states), "alphabet": str(self.alphabet), "start": str(self.start), "tf": str(self.tf), "accepting": str(self.accs)}
        jsonString = json.dumps(DFA_dict)
        
        return jsonString

    # Breadth First Search Counter-Exampling method
    def BFS_ce(self, h_DFA):
        visited, queue = {}, [["", self.start]]

        if self.label("") != h_DFA.label(""):
            return ""

        while queue:
            vertex_lst = queue.pop(0)
            
            if vertex_lst[1] not in visited:
                visited[vertex_lst[1]] = 1
            else:
                visited[vertex_lst[1]] += 1


            if visited[vertex_lst[1]] <= 2:

                for sym, next_state in self.tf[vertex_lst[1]].items():
    
                    if self.label(vertex_lst[0] + sym) != h_DFA.label(vertex_lst[0] + sym):
                        return vertex_lst[0] + sym

                    queue.append([vertex_lst[0] + sym, next_state])
        
        return "None"