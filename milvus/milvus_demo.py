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

from pymilvus import (
    connections,
    utility,
    FieldSchema,
    CollectionSchema,
    DataType,
    Collection,
)


connections.connect(host="127.0.0.1", port="19530")

# create collection

if utility.has_collection("liu_lang_di_qiu"):
    utility.drop_collection("liu_lang_di_qiu")

fields = [
    FieldSchema(name="id", dtype=DataType.INT64, is_primary=True),
    FieldSchema(
        name="embedding", dtype=DataType.FLOAT_VECTOR, dim=sentence_embeddings.shape[1]
    ),
]
schema = CollectionSchema(fields=fields)
collection = Collection(name="liu_lang_di_qiu", schema=schema)

entities = [
    [i for i in range(sentence_embeddings.shape[0])],
    [x for x in sentence_embeddings],
]
collection.insert(entities)
collection.flush()

# create index

index_params = {"index_type": "IVF_FLAT", "metric_type": "L2", "params": {"nlist": 10}}
collection.create_index(
    field_name="embedding", index_params=index_params, index_name="embedding_index"
)

# search

collection = Collection("liu_lang_di_qiu")
collection.load()

search_params = {
    "params": {"nprobe": 3},
}
vectors_to_search = model.encode(["太阳炸了"])

# result

result = collection.search(
    vectors_to_search,
    "embedding",
    search_params,
    limit=3,
    output_fields=["id"],
)

for id in result[0].ids:
    print(df["sentence"].iloc[id])

# benchmark

import time


costs = []
for x in range(10000):
    t0 = time.time()
    result = collection.search(
        vectors_to_search,
        "embedding",
        search_params,
        limit=3,
        output_fields=["id"],
    )
    t1 = time.time()
    costs.append(t1 - t0)
print("平均耗时 %7.3f ms" % ((sum(costs) / len(costs)) * 1000.0))
