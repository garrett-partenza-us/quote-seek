from sklearn.manifold import TSNE
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd

# Load the pre-saved embeddings
df = pd.read_pickle("./vectors.pkl")

# Extract fine-tuned and base embeddings
fine_tuned_embeddings = np.array([x.flatten() for x in df.fine_tuned_embedding])
base_embeddings = np.array([x.flatten() for x in df.base_embedding])

# Initialize t-SNE with 2 components
tsne_fine_tuned = TSNE(n_components=2, random_state=42)
tsne_base = TSNE(n_components=2, random_state=42)

# Apply t-SNE transformation to embeddings
X_tsne_fine_tuned = tsne_fine_tuned.fit_transform(fine_tuned_embeddings)
X_tsne_base = tsne_base.fit_transform(base_embeddings)

# Create subplots for both fine-tuned and base embeddings
fig, axes = plt.subplots(1, 2, figsize=(16, 6))

# Plot Fine-Tuned Embeddings
axes[0].scatter(X_tsne_fine_tuned[:, 0], X_tsne_fine_tuned[:, 1], c='red', alpha=0.7, edgecolor='black', s=40)
axes[0].set_xlabel('t-SNE Component 1')
axes[0].set_ylabel('t-SNE Component 2')
axes[0].set_title('2D t-SNE: Fine-Tuned RoBERTa Embeddings')

# Plot Base Embeddings
axes[1].scatter(X_tsne_base[:, 0], X_tsne_base[:, 1], c='blue', alpha=0.7, edgecolor='black', s=40)
axes[1].set_xlabel('t-SNE Component 1')
axes[1].set_ylabel('t-SNE Component 2')
axes[1].set_title('2D t-SNE: Base RoBERTa Embeddings')

# Adjust layout and show the plots
plt.tight_layout()
plt.show()
