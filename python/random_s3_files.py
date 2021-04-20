import random

import boto3
import lorem
from botocore.exceptions import ClientError
import json


def get_s3_client():
    """Instantiates a new S3 client using boto3"""
    return boto3.Session(profile_name="apt304").client("s3")


def random_document(doc_id):
    return {'size': random.randint(1, 101),
            'duration': random.randint(1, 101),
            'text': lorem.paragraph(),
            'id': doc_id,
            }


def add_documents_to_bucket(s3_client, root, key, num_of_docs):
    for i in range(num_of_docs):
        doc = random_document(i)
        s3_client.put_object(
            Body=json.dumps(doc),
            Bucket=root,
            Key="{}-{}.json".format(key, i, num_of_docs))


def populate():
    root = "apt304-spencer-test"
    buckets = ["argentina-legislation", "brazil-legislation", "omni-search"]
    s3_client = get_s3_client()
    for bucket in buckets:
        add_documents_to_bucket(s3_client, root, "fn-s3-es-test/dev/{}/{}".format(bucket, bucket), 1000)


if __name__ == "__main__":
    populate()
