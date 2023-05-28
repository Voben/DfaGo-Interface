# custom made DFA class 
from .dfa import DFA, ToDfaGoJSON
import socket

# this function will create a Lstar object
# and start learning
def LearnLstar(ip, port, buffer_size):
    lstar_obj = Lstar()
    ip_str = ip.decode()

    lstar_obj.init(ip_str, port, buffer_size)
    jsonDFA = lstar_obj.learn_DFA(["0", "1"])

    return jsonDFA.encode()

# Lstar class
class Lstar:
    # separate init function
    def init(self, ip, port, buffer_size):
        self.buffer_size = buffer_size
    
        self.s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.s.connect((ip, int(port)))

    # ask a membership query to the connected TCP server
    def ask_mq(self, word):
        self.num_mq += 1

        query = "mq " + word 

        self.s.sendall(query.encode())
        data = self.s.recv(self.buffer_size).decode()
        if data == "true":
            return True
        else:
            return False 

    # ask equivalence query to the connected TCP server
    def ask_eq(self, json_dfa):
        self.num_eq += 1
        query = "eq " + json_dfa

        self.s.sendall(query.encode())
        data = self.s.recv(self.buffer_size)

        dataStr = data.decode()
        if dataStr == "empty":
            dataStr = ""

        return dataStr

    # sends end message to TCP server to close the connection
    def end_ws(self):
        self.s.send("end".encode())

    # will learn a DFA based on the examples and regular expression provided
    def learn_DFA(self, alphabet = ['0', '1']):
        self.num_mq = 0
        self.num_eq = 0

        self.alphabet = alphabet
        self.iter = 0
        
        qQ = [""]
        tT = [""]

        while True:
            # expand states and transition function 
            qQ, tf = self.close(qQ, tT)

            # get final states in q
            final_states = []
            for pos in qQ:
                if self.ask_mq(pos):
                    final_states.append(pos)

            # create hypothesis dfa
            self.hDFA = DFA(qQ, self.alphabet, "", tf, final_states)

            # determine if dfa properly classifies the examples provided
            res = self.ask_eq(ToDfaGoJSON(self.hDFA))
            # if so return the hypothesis dfa
            if res == "None": # if res True then hDFA classifies L (or the examples provided)

                self.end_ws()
                return ToDfaGoJSON(self.hDFA)
            
            # consider the counter example and create a new state and test word to be considered
            new_state, new_test_word = self.add_test_word(qQ, tf, res) # res should contain a counterexample that H does not classify
            
            if new_state not in qQ:
                qQ.append(new_state)

            if new_test_word not in tT:
                tT.append(new_test_word)

            self.iter += 1


    # function to check for t equivalence
    def are_indistinguishable(self, T, u, v):
        # go through the suffixes stored in T
        for t in T:
            if self.ask_mq( (u+t) ) != self.ask_mq( (v+t)):
                return False

        return True

    # main loop sequence in l star
    def close(self, qQ, tT):
        # create a new empty tf, i to go through all states, create a copy of states to avoid infinite loop
        tf = {}
        i = 0
        newQ = qQ.copy()

        # go through all states
        while i < len(qQ):
            
            # go through the alphabet
            for a in self.alphabet:
                
                # get a new state
                q = qQ[i]

                # go through all states again
                for r in qQ:

                    # check if a state and a character in the alphabet lead to another state meaning that these two states are T-equivalent
                    if self.are_indistinguishable(tT, (q + a), r):
                        
                        # avoid key error if this is the first instance of q in tf
                        if q not in tf:
                            tf[q] = {}

                        # link the two states in the tf
                        tf[q][a] = r    
                        break
                
                # avoid key error if first instance of q in tf
                if q not in tf:
                    tf[q] = {}

                # if a not in the tf[q] append the state and letter as a new state and create a new entry in tf
                if a not in tf[q]:
                    newQ.append(q + a)
                    tf[q][a] = q + a

            i += 1
        return newQ, tf

    # used to treat a new counter example correctly
    # will go through the transition function accordingly depending on the counter example provided
    def add_test_word(self, Q, tf, ce):
        # q will store the string as we go through it and i will store the index
        q = ""
        i = 0

        # infinite loop
        while True:
            # try except that will trigger if the state is not yet in tf, same thing is carried out if triggered
            # check if the transition from the current state and the rest of string leads to the counter example
            try:
                if self.ask_mq( (tf[q][ce[i]] + ce[i+1:]) ) != self.ask_mq(ce):
                    return (q + ce[i]), ce[i+1:]

            except KeyError:
                return (q + ce[i]), ce[i+1:]                

            # if not store q as the current part of the string and update the i
            q = tf[q][ce[i]]
            i += 1