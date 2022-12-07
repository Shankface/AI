from sklearn.model_selection import train_test_split
import pandas as pd
import random

df = pd.read_csv("Heart Disease Dataset.csv")

data = df.sample(frac = 1)

feature_cols = ['age',	'sex',	'cp',	'trestbps',	'chol',	'fbs',	'restecg',	'thalach',	'exang',	'oldpeak',	'slope',	'ca',	'thal']
X = data[feature_cols]
Y = data['target']

X_train, X_test, Y_train, Y_test = train_test_split(X, Y, test_size = 0.25)
# X_train = (X_train - X_train.min())/(X_train.max()-X_train.min())
mean = X_train.mean()i thn
std = X_train.std()

X_train = (X_train - mean)/(std)
X_train = X_train.to_numpy()
Y_train = Y_train.to_numpy()
# X_test = (X_test - X_train.min())/(X_test.max()-X_train.min())
X_test = (X_test - mean)/(std)
X_test = X_test.to_numpy()
Y_test = Y_test.to_numpy()

f = open("heart_train_data.txt", "w")
f.write(str(len(X_train)) + " " + str(len(X_train[0])) + " 1\n")
for i,row in enumerate(X_train):
    f.write(" ".join(str(f'{val:.3f}') for val in row))
    f.write(" " + str(Y_train[i]) + "\n")

f = open("heart_test_data.txt", "w")
f.write(str(len(X_test)) + " " + str(len(X_test[0])) + " 1\n")
for i,row in enumerate(X_test):
    f.write(" ".join(str(f'{val:.3f}') for val in row))
    f.write(" " + str(Y_test[i]) + "\n")

for h in range(5,20):
    f = open("heart_init_" + str(h) + ".txt", "w")
    f.write(str(len(X_train[0])) + " " + str(h) + " 1\n")
    for i in range(h):
        for j in range(len(X_train[0]) + 1):
            f.write(str(f'{random.uniform(0, 1):.3f}') + " ")
        f.write("\n")

    for i in range(1):
        for j in range(h + 1):
            f.write(str(f'{random.uniform(0, 1):.3f}') + " ")
        f.write("\n")

