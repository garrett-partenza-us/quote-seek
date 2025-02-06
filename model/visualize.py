import torch
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from sklearn.decomposition import PCA
from encoder import Encoder
from transformers import BertForSequenceClassification, BertModel, BertTokenizer

# Load the fine-tuned and base models
device = torch.device("mps" if torch.backends.mps.is_available() else "cpu")

fine_tuned_model = BertForSequenceClassification.from_pretrained("./fine_tuned_bert")
base_model = BertModel.from_pretrained("bert-base-uncased")
tokenizer = BertTokenizer.from_pretrained("./fine_tuned_bert")

fine_tuned_model.to(device)
base_model.to(device)

# Initialize the encoder
encoder = Encoder("bert-base-uncased")

# Load the pre-saved embeddings
df = pd.read_pickle("./vectors.pkl")

# Load PCA components
fine_tuned_embeddings = np.array(list(x.flatten() for x in df.fine_tuned_embedding))
base_embeddings = np.array(list(x.flatten() for x in df.base_embedding))

# Fit PCA models
pca_fine_tuned = PCA(n_components=2)
pca_base = PCA(n_components=2)

X_pca_fine_tuned = pca_fine_tuned.fit_transform(fine_tuned_embeddings)
X_pca_base = pca_base.fit_transform(base_embeddings)

def plot_embeddings_with_query(query_text, plot_query=True):
    # Encode the query text
    encodings, _ = encoder.encode_query(query_text)
    input_ids = encodings['input_ids'].to(device)
    attention_mask = encodings['attention_mask'].to(device)

    with torch.no_grad():
        fine_tuned_outputs = fine_tuned_model.bert(input_ids=input_ids, attention_mask=attention_mask)
        fine_tuned_hidden_states = fine_tuned_outputs.last_hidden_state
        fine_tuned_vector = fine_tuned_hidden_states[:, 0, :].detach().cpu().numpy()

        base_outputs = base_model(input_ids=input_ids, attention_mask=attention_mask)
        base_hidden_states = base_outputs.last_hidden_state
        base_vector = base_hidden_states[:, 0, :].detach().cpu().numpy()

    # Project the query vector to PCA space
    query_fine_tuned_pca = pca_fine_tuned.transform(fine_tuned_vector)
    query_base_pca = pca_base.transform(base_vector)

    # Plot both PCA graphs
    fig, axes = plt.subplots(1, 2, figsize=(16, 6))

    # Plot Fine-Tuned PCA graph
    axes[0].scatter(X_pca_fine_tuned[:, 0], X_pca_fine_tuned[:, 1], c='red', alpha=0.7, edgecolor='black', s=40)
    if plot_query:
        axes[0].scatter(query_fine_tuned_pca[0, 0], query_fine_tuned_pca[0, 1], c='green', marker='x', s=100, label='Query')
    axes[0].set_xlabel('Principal Component 1')
    axes[0].set_ylabel('Principal Component 2')
    axes[0].set_title('2D PCA: Fine-Tuned RoBERTa Embeddings')
    axes[0].legend()

    # Plot Base PCA graph
    axes[1].scatter(X_pca_base[:, 0], X_pca_base[:, 1], c='blue', alpha=0.7, edgecolor='black', s=40)
    if plot_query:
        axes[1].scatter(query_base_pca[0, 0], query_base_pca[0, 1], c='green', marker='x', s=100, label='Query')
    axes[1].set_xlabel('Principal Component 1')
    axes[1].set_ylabel('Principal Component 2')
    axes[1].set_title('2D PCA: Base RoBERTa Embeddings')
    axes[1].legend()

    plt.tight_layout()
    plt.show()

if __name__ == "__main__":
    # Prompt the user for a text query
    query_text = input("Enter the query text: ")
    plot_embeddings_with_query(query_text)
