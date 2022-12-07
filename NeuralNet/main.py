from random import random
import math
from tqdm import tqdm
import sys
import os

# Function to parse a network file and passes it to network creator function
def process_network_file(path):
    init_hidden_w, init_output_w = [],[]
    f = open(path, "r")
    n_inp,n_hidden,n_outputs = f.readline().split()

    for i in range(int(n_hidden)):
        init_hidden_w.append(f.readline().split())
    for i in range(int(n_outputs)):
        init_output_w.append(f.readline().split())

    return initialize_network(init_hidden_w, init_output_w)
    
# Function to parse train or test data file
def process_data_file(path):
    data = []
    f = open(path, "r")
    row = f.readline()
    while row:
        row = f.readline().split()
        if row:
            data.append(list(map(float, row)))
    return data

# Function to create a Neural Network with a hidden layer and an output layer
def initialize_network(init_hidden_w, init_output_w):
    network = list()
    hidden_layer = [{'weights':[float(w) for w in layer]} for layer in init_hidden_w]
    network.append(hidden_layer)
    output_layer = [{'weights':[float(w) for w in layer]} for layer in init_output_w]
    network.append(output_layer)
    return network

# Sigmoid activation function
def activation(weights, input):
    act = -weights[0]
    for i in range(1,len(weights)):
        act += weights[i] * input[i-1]
    return 1.0 / (1.0 + math.exp(-act))

# Derivative of sigmoid activation function
def activation_deriv(output):
    return output*(1.0-output)

# Forward Propogation
def forward_prop(network, data_row):
    input = data_row

    for layer in network:
        new_inp = []
        for node in layer:
            node["output"] = activation(node["weights"], input)
            new_inp.append(node["output"])
        input = new_inp
    return input

# Back Propogation
def back_prop(network, data_row, num_outputs):
    expected = data_row[-num_outputs:]
    inputs = data_row[:-num_outputs]

    for j in range(num_outputs):
        node = network[1][j]
        node['delta'] = (expected[j] - node['output']) * activation_deriv(node['output'])
    
    for i in range(len(network[0])):
        sum = 0
        for j in range(num_outputs):
            node = network[1][j]
            sum += node["weights"][i+1] * node["delta"]

        network[0][i]["delta"] =  activation_deriv(activation(network[0][i]["weights"], inputs)) * sum

# Function to update weights
def update_weights(network, data_row, num_outputs, lr):
    for i,layer in enumerate(network):
        inputs = data_row[:-num_outputs]

        if i == 1:
            inputs = [node['output'] for node in network[0]]

        for node in layer:
            for j in range(len(inputs)):
                node['weights'][j+1] += lr * node['delta'] * inputs[j]
            node['weights'][0] += lr * node['delta'] * (-1)

# Function that takes network and trains it using inputted data and outputs it to a file
def train_network(network, data, lr, epochs, outfile):
    num_outputs = len(network[1])
    print("\nTraining for " + str(epochs) + " epochs with a learning rate of " + str(lr))
    print("--------------------------------")
    for epoch in range(epochs):
        error_sum = 0
        for data_row in data:
            expected = data_row[-num_outputs:]
            outputs = forward_prop(network, data_row)
            error_sum += sum([(expected[i]-outputs[i])**2 for i in range(len(expected))])
            back_prop(network, data_row, num_outputs)
            update_weights(network, data_row, num_outputs, lr)
        print("Epoch: %d/%d, error = %s" % (epoch+1, epochs, str(f'{error_sum:.3f}')))
    print("--------------------------------\nTraining Complete")

    f = open(outfile, "w")
    f.write(str(len(network[0][0]["weights"]) - 1) + " " + str(len(network[0])) + " " + str(num_outputs) + "\n")
    for layer in network:
        for node in layer:
            f.write(" ".join(str(f'{w:.3f}') for w in node["weights"]) + "\n")

# Function that takes network and tests it using inputted data and outputs the results to a file
def test_network(network, data, outfile):
    num_outputs = len(network[1])
    m = {"A": 0, "B": 0, "C": 0, "D": 0, "acc": 0, "prec": 0, "recall": 0, "F1": 0}
    test_metrics = []
    A,B,C,D,acc_macro,prec_macro,recall_macro,F1_macro = 0,0,0,0,0,0,0,0

    for i in range(num_outputs):
        test_metrics.append(m.copy())

    print("\nTesting...")
    for data_row in data:
        actual = data_row[-num_outputs:]
        pred = forward_prop(network, data_row)

        for i in range(num_outputs):
            if ((pred[i] >= 0.5) and (actual[i] == 1)):
                test_metrics[i]["A"] += 1
            elif ((pred[i] >= 0.5) and (actual[i] == 0)):
                test_metrics[i]["B"] += 1
            elif ((pred[i] < 0.5) and (actual[i] == 1)):
                test_metrics[i]["C"] += 1
            elif ((pred[i] < 0.5) and (actual[i] == 0)):
                test_metrics[i]["D"] += 1

    for i in range(num_outputs):
        test_metrics[i]["acc"] = (test_metrics[i]["A"]+test_metrics[i]["D"])/(test_metrics[i]["A"]+test_metrics[i]["B"]+test_metrics[i]["C"]+test_metrics[i]["D"])
        test_metrics[i]["prec"] = test_metrics[i]["A"]/(test_metrics[i]["A"]+test_metrics[i]["B"])
        test_metrics[i]["recall"] = test_metrics[i]["A"]/(test_metrics[i]["A"]+test_metrics[i]["C"])
        test_metrics[i]["F1"] = (2 * test_metrics[i]["prec"] *  test_metrics[i]["recall"])/(test_metrics[i]["prec"] + test_metrics[i]["recall"])
        
        A += test_metrics[i]["A"]
        B += test_metrics[i]["B"]
        C += test_metrics[i]["C"]
        D += test_metrics[i]["D"]
        
        acc_macro += test_metrics[i]["acc"]
        prec_macro += test_metrics[i]["prec"]
        recall_macro += test_metrics[i]["recall"]

    acc_micro = (A+D)/(A+B+C+D)
    prec_micro = A/(A+B)
    recall_micro = A/(A+C)
    F1_micro = (2 * prec_micro * recall_micro)/(prec_micro + recall_micro)

    acc_macro = acc_macro/num_outputs
    prec_macro = prec_macro/num_outputs
    recall_macro = recall_macro/num_outputs
    F1_macro = (2 * prec_macro * recall_macro)/(prec_macro + recall_macro)

    print("Testing Complete")
    f = open(outfile, "w")
    for categ in test_metrics:
        for i,metric in enumerate(categ):
            if i < 4:
                f.write(str(categ[metric]) + " ")
            elif i < 7:
                f.write(str(f'{categ[metric]:.3f}') + " ")
            else:
                f.write(str(f'{categ[metric]:.3f}') + "\n")
    f.write(str(f'{acc_micro:.3f}') + " " + str(f'{prec_micro:.3f}') + " " + str(f'{recall_micro:.3f}') + " " + str(f'{F1_micro:.3f}' + '\n'))
    f.write(str(f'{acc_macro:.3f}') + " " + str(f'{prec_macro:.3f}') + " " + str(f'{recall_macro:.3f}') + " " + str(f'{F1_macro:.3f}' + '\n'))

# Main function
if __name__ == "__main__":
    typ = input("Input 'train' or 'test': ")

    if typ == 'train':
        init_weights_path = input("Input the untrained network file: ")
        train_data_path = input("Input training data file: ")
        output_path = input("Input the output file for the trained network: ")
        epochs = int(input("Input amount of epochs to train for: "))
        lr = float(input("Input learning rate: "))

        network = process_network_file(init_weights_path)
        data = process_data_file(train_data_path)
        train_network(network, data, lr, epochs, output_path)

    if typ == 'test':
        trained_weights_path = input("Input the trained network file: ")
        test_data_path = input("Input testing data file: ")
        output_path = input("Input the output file for testing results: ")

        network = process_network_file(trained_weights_path)
        data = process_data_file(test_data_path)
        test_network(network, data, output_path)