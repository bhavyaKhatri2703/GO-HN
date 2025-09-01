import psycopg2
import numpy as np

def connect_db():
    conn = psycopg2.connect(
            host="localhost",
            port  = 5555,
            database="hackernews",
            user="postgres",
            password="postgres"
        )
    return conn


def createTables(conn):
    cursor = conn.cursor()

    cursor.execute("""
        CREATE TABLE IF NOT EXISTS newStories (
            id SERIAL PRIMARY KEY,
            by TEXT,
            type TEXT,
            text TEXT,
            url TEXT,
            title TEXT,
            full_text TEXT,
            score INT,
            sem_embedding vector(384),
            bm25_embedding bm25vector
            )
        """)


    cursor.execute("""
        CREATE TABLE IF NOT EXISTS topStories (
            id SERIAL PRIMARY KEY,
            by TEXT,
            type TEXT,
            text TEXT,
            url TEXT,
            title TEXT,
            full_text TEXT,
            score INT,
            sem_embedding vector(384),
            bm25_embedding bm25vector
        )
        """)

    conn.commit()
    cursor.close()

    print("Tables created/verified")


def saveNewStory(conn,story):
    cursor = conn.cursor()

    embedding = story['embeddings']
    if isinstance(embedding, np.ndarray):
        embedding = embedding.tolist()

    cursor.execute("""
                  INSERT INTO topStories (id,by,type,text,url,title,full_text,score,sem_embedding)
                  VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s)
                  ON CONFLICT (id) DO NOTHING
              """,(story['id'],story['by'],story['type'],story['text'],story['url'],story['title'],story['full_text'],story['score'],embedding)
          )
    print(f"saved : {story['id']}")
    conn.commit()
    cursor.close()

def saveTopStory(conn,story):
    cursor = conn.cursor()

    embedding = story['embeddings']
    if isinstance(embedding, np.ndarray):
        embedding = embedding.tolist()

    cursor.execute("""
                INSERT INTO topStories (id,by,type,text,url,title,full_text,score,sem_embedding)
                VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s)
                ON CONFLICT (id) DO NOTHING
            """,(story['id'],story['by'],story['type'],story['text'],story['url'],story['title'],story['full_text'],story['score'],embedding)
        )

    print(f"saved : {story['id']}")

    conn.commit()
    cursor.close()

def deleteNewStory(conn , id):
    cursor = conn.cursor()

    cursor.execute("""
              DELETE FROM newStories WHERE id = %s
          """, (id,))

    conn.commit()
    cursor.close()

def deleteTopStory(conn , id):
    cursor = conn.cursor()

    cursor.execute("""
            DELETE FROM topStories WHERE id = %s
        """, (id,))

    conn.commit()
    cursor.close()
