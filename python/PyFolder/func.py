import sys
import requests
import json

# functions which convert variables to different types
def StringToBytes(str):
    return str.encode("utf-8")

def StringToByteArray(str):
    return bytearray(str, "utf-8")

def BytesToString(by):
    return by.decode()

def ByteArrayToString(byAr):
    return str(byAr, "utf-8")

# get the refCount of a provided variable
def refCount(obj):
    return sys.getrefcount(obj)

# print a Python variable from Python
def printPyObject(obj):
    print(obj)

# a function that provided a DfaGo dataset 
# as a Python bytes that will convert
# it to a Python variable
def getPySet(sample_lst):
    set_lst = sample_lst.decode().split(",")
    pos_lst = []
    neg_lst = []
    idx = 0
    
    for i in range(0, len(set_lst)):
        if set_lst[i] == 'positive':
            pass
        elif set_lst[i] == 'negative':
            idx = i
            break
        else:
            pos_lst.append(set_lst[i])

    for i in range(idx, len(set_lst)):
        if set_lst[i] == "negative":
            continue
        else:
            neg_lst.append(set_lst[i])

    ret_set = {"positive": pos_lst, "negative":neg_lst}

    return ret_set

