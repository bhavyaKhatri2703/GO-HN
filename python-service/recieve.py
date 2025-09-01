import requests
import pika
import json
import threading
from queries import connect_db,createTables,saveTopStory,saveNewStory,deleteTopStory,deleteNewStory
from embeddings import generateEmbeddings


conn = connect_db()
createTables(conn)


def getStories(id) :
    url = f"https://hacker-news.firebaseio.com/v0/item/{id}.json"

    response = requests.get(url)

    data = response.json()

    story_data = {
        'by' : data.get('by',''),
        'id' : data.get('id',0),
        'type' : data.get('type' , ''),
        'text' : data.get('text', ""),
        'url' : data.get('url' , ""),
        'title' : data.get('title',""),
        'full_text' : data.get('title' , "") + ' ' + data.get('text' , ""),
        'score' : data.get('score' , 0),
        'embeddings' : generateEmbeddings(data.get('title' , "") + " " + data.get('text' , ""))
        }
    return story_data


def add_callback(conn , queue_name):
    def callback(ch , method , properties , body):
        try:
            message_data = json.loads(body.decode('utf-8'))
            story_type = message_data.get("type")
            ids = message_data.get("ids")

            if story_type == 'new':
                for id in ids :
                    story_data = getStories(id)
                    saveNewStory(conn,story_data)
            else:
                for id in ids :
                    story_data = getStories(id)
                    saveTopStory(conn,story_data)


        except json.JSONDecodeError as e:
            print(f"[{queue_name}] Failed to decode message: {e}")

    return callback

def delete_callback(conn,queue_name):
    def callback(ch , method , properties , body):
        try:
             message_data = json.loads(body.decode('utf-8'))
             story_type = message_data.get("type")
             ids = message_data.get("ids")

             if story_type == 'new':
                 for id in ids :
                     deleteNewStory(conn,id)

             else:
                 for id in ids :
                     story_data = getStories(id)
                     deleteTopStory(conn,story_data)

        except json.JSONDecodeError as e:
            print(f"[{queue_name}] Failed to decode message: {e}")

    return callback

def consume_queue(queue_name , callback_func):
    connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
    channel = connection.channel()


    channel.queue_declare(queue=queue_name, durable=True)

    channel.basic_consume(queue=queue_name, on_message_callback=callback_func, auto_ack=True)

    print(f"[{queue_name}] Waiting for messages...")
    channel.start_consuming()


queues = ['hn_add', 'hn_delete']



queue_callbacks = {
    'hn_add': add_callback(conn, 'hn_add'),
    'hn_delete': delete_callback(conn, 'hn_delete'),
}

for queue_name, cb in queue_callbacks.items():
    t = threading.Thread(target=consume_queue, args=(queue_name, cb), daemon=True)
    t.start()

try:
    while True:
        pass
except KeyboardInterrupt:
    print("Shutting down.")
