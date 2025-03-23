import os
from tqdm import tqdm
import pandas as pd
from openai import OpenAI
from dotenv import load_dotenv
from concurrent.futures import ThreadPoolExecutor

# Load environment variables
load_dotenv()

# Generate OpenAI embeddings
def generate_openai_embeddings(meditations):
    # Function to fetch embeddings for a single quote
    def get_embedding(quote):
        client = OpenAI(api_key=os.getenv("API_KEY"))
        response = client.embeddings.create(
            input=quote,
            model="text-embedding-3-small"
        )
        return response.data[0].embedding

    # Initialize the list to store results
    results = []

    # Use ThreadPoolExecutor for parallelization
    with ThreadPoolExecutor() as executor:
        # Submit tasks and ensure results stay in order by indexing the output
        futures = {executor.submit(get_embedding, quote): idx for idx, quote in
                   enumerate(meditations["chunk"])}
        
        # Retrieve results as they complete and maintain the original order
        for future in tqdm(futures, desc="Fetching embeddings"):
            idx = futures[future]
            result = future.result()  # Get the result from the future
            results.append((idx, result))  # Store the index to maintain order

    # Sort results by the original index and extract the embeddings in order
    results.sort(key=lambda x: x[0])
    embeddings = [embedding for _, embedding in results]

    # Add the embeddings to the dataframe
    meditations["openai_embedding"] = embeddings

    return meditations


if __name__ == '__main__':
    # Ensure the MEDITATIONS environment variable is loaded and use it to read CSV
    meditations = pd.read_csv(os.getenv("MEDITATIONS"))
    df_with_openai_embeddings = generate_openai_embeddings(meditations)
    df_with_openai_embeddings.to_csv("meditations.csv", index=False)
