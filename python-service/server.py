# server.py

import grpc
from concurrent import futures
from proto import embeddings_pb2
from proto import embeddings_pb2_grpc
from embeddings import generateEmbeddings

class EmbeddingsService(embeddings_pb2_grpc.EmbeddingsServiceServicer):
    def GetEmbeddings(self, request, context):
        interests = request.interests

        text = " ".join(interests)
        embeddings = generateEmbeddings(text)
        return embeddings_pb2.InterestsResponse(embeddings=embeddings)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    embeddings_pb2_grpc.add_EmbeddingsServiceServicer_to_server(EmbeddingsService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server started on port 50051")
    server.wait_for_termination()

if __name__ == "__main__":
    serve()
