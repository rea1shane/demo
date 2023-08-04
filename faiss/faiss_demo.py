# Doc: https://soulteary.com/2022/09/03/vector-database-guide-talk-about-the-similarity-retrieval-technology-from-metaverse-big-company-faiss.html

import pandas as pd
import os


df = pd.read_csv(
    os.path.dirname(__file__) + "/../_data/liu-lang-di-qiu.txt",
    sep="#",
    header=None,
    names=["sentence"],
)
# print(df)

####################

from sentence_transformers import SentenceTransformer


model = SentenceTransformer("uer/sbert-base-chinese-nli")
sentences = df["sentence"].tolist()
sentence_embeddings = model.encode(sentences)
# print(sentence_embeddings.shape)

####################

import faiss


dimension = sentence_embeddings.shape[1]

# index = faiss.IndexFlatL2(dimension)
# index.add(sentence_embeddings)

quantizer = faiss.IndexFlatL2(dimension)
nlist = 10
index = faiss.IndexIVFFlat(quantizer, dimension, nlist)
index.train(sentence_embeddings)
index.add(sentence_embeddings)
index.nprobe = 3

# print(index.ntotal)

search = model.encode(["太阳炸了"])
D, I = index.search(search, 3)
print(df["sentence"].iloc[I[0]])

####################

import time


costs = []
for x in range(10000):
    t0 = time.time()
    D, I = index.search(search, 3)
    t1 = time.time()
    costs.append(t1 - t0)
print("平均耗时 %7.3f ms" % ((sum(costs) / len(costs)) * 1000.0))

# 0.052 ms
# 0.018 ms
