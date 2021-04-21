import multiprocessing
import queue
import random
import time

import boto3
import lorem
from botocore.exceptions import ClientError
import json


def get_aws_session():
    """Instantiates a new S3 client using boto3"""
    return boto3.Session(profile_name="apt304")


def random_text(size):
    text = ""
    for i in range(size):
        text += lorem.paragraph()
    return text


def random_document(index, doc_id):
    return {'size': random.randint(1, 101),
            'duration': random.randint(1, 101),
            'text': random_text(random.randint(8, 50)),
            'id': doc_id,
            'index': index}


def document_generator(indexes, document_count):
    for index in indexes:
        for i in range(document_count):
            yield random_document(index, i)


def process(_document_queue, completed_event, counter, root):
    session = get_aws_session()
    s3 = session.client("s3")

    while not _document_queue.empty() or not completed_event.is_set():
        document = None
        try:
            document = _document_queue.get(False)
            counter.value += 1
        except queue.Empty:
            pass

        if document:
            obj = s3.put_object(
                Body=json.dumps(document),
                Bucket=root,
                Key="{}/{}/{}-{}.json".format(root, document["index"], document["index"], document["id"]))


def status(start, since, generated, newly_gen, uploaded, newly_upl, total):
    now = time.time()
    print('Generated {}/{} ({} per s). Uploaded {}/{} ({} per s). Elapsed time {}'.format(generated, total,
                                                                                          newly_gen / (now - since),
                                                                                          uploaded, total,
                                                                                          newly_upl / (now - since),
                                                                                          now - start))

def populate():
    ROOT = "apt304-spencer-test"
    INDEXES = ["argentina_legislation", "australia_legislation", "bill", "brazil_legislation", "canada_legislation",
               "twitter"]
    DOCUMENT_COUNT = 10000

    # via processes
    PROCCESS_COUNT = 50
    document_queue = multiprocessing.Queue()
    completed_event = multiprocessing.Event()
    counter = multiprocessing.Value('i', 0)
    processes = []

    for i in range(PROCCESS_COUNT):
        p = multiprocessing.Process(target=process, args=(document_queue, completed_event, counter, ROOT))
        p.start()
        processes.append(p)

    start, since = time.time(), time.time()
    newly_gen, total_gen, prev_upl = 0, 0, 0
    for doc in document_generator(INDEXES, DOCUMENT_COUNT):
        document_queue.put(doc)
        newly_gen, total_gen = newly_gen + 1, total_gen + 1

        if total_gen % 1000 == 0:
            newly_upl = counter.value - prev_upl
            prev_upl = counter.value
            status(start, since, total_gen, newly_gen, prev_upl, newly_upl, len(INDEXES) * DOCUMENT_COUNT)
            since, newly_gen = time.time(), 0

    completed_event.set()
    since = time.time()
    newly_gen = 0

    while not document_queue.empty() or not completed_event.is_set():
        newly_upl = counter.value - prev_upl
        prev_upl = counter.value
        status(start, since, total_gen, newly_gen, prev_upl, newly_upl, len(INDEXES) * DOCUMENT_COUNT)
        since, newly_gen = time.time(), 0
        time.sleep(5)

    for p in processes:
        p.join()


if __name__ == "__main__":
    populate()
