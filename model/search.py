import torch
from encoder import Encoder
import numpy as np
import pandas as pd
from tqdm import tqdm
import json
from sklearn.metrics.pairwise import cosine_similarity
from transformers import BertTokenizer, BertModel
from transformers import BertForSequenceClassification, BertTokenizer

model_name = "bert-base-uncased"
encoder = Encoder(model_name)

fine_tuned_model = BertForSequenceClassification.from_pretrained("./fine_tuned_bert")
base_model = BertModel.from_pretrained(model_name)

device = torch.device("mps" if torch.backends.mps.is_available() else "cpu")
fine_tuned_model.to(device)
base_model.to(device)

df = pd.read_pickle('./vectors.pkl')
fine_tuned_database = np.array([x.flatten() for x in df.fine_tuned_embedding])
base_database = np.array([x.flatten() for x in df.base_embedding])

if __name__ == '__main__':

    while True:

        text = input("\nPlease enter your text: ")
        choice = None
        while True:
            choice = input("Enter 0 for base or 1 for fine-tuned: ")
            if choice in ['0', '1']:
                model = fine_tuned_model if choice == '1' else base_model
                break

        encodings, sentence_ids = encoder.encode_query(text)
        input_ids = encodings['input_ids']
        attention_mask = encodings['attention_mask']
        texts = [encoder.lookup_index(x) for x in sentence_ids]
        """
        with torch.no_grad():
            if choice == '1':
                output = fine_tuned_model.bert(input_ids=input_ids.to(device),
                                                  attention_mask=attention_mask.to(device))
                hidden = output.last_hidden_state
                vector = hidden[:, 0, :].detach().cpu().flatten()
                database = fine_tuned_database
            else:
                output = base_model(input_ids=input_ids.to(device),
                                    attention_mask=attention_mask.to(device))
                hidden = output.last_hidden_state
                vector = hidden[:, 0, :].detach().cpu().flatten()
                database = base_database
            """
        with torch.no_grad():
            if choice == '1':
                output = fine_tuned_model.bert(input_ids=input_ids.to(device),
                                               attention_mask=attention_mask.to(device))
                hidden = output.last_hidden_state.cpu()
                # Take the average of all token embeddings, ignoring padding tokens
                attention_mask_expanded = attention_mask.unsqueeze(-1).expand(hidden.size())  # To apply attention mask
                sum_hidden = torch.sum(hidden * attention_mask_expanded, dim=1)  # Sum over all tokens
                count_non_padding = attention_mask.sum(dim=1, keepdim=True)  # Count of non-padding tokens
                vector = sum_hidden / count_non_padding  # Average embedding
                vector = vector.detach().cpu().flatten()  # Flatten the vector
                database = fine_tuned_database
            else:
                output = base_model(input_ids=input_ids.to(device),
                                    attention_mask=attention_mask.to(device))
                hidden = output.last_hidden_state
                # Take the average of all token embeddings, ignoring padding tokens
                attention_mask_expanded = attention_mask.unsqueeze(-1).expand(hidden.size())  # To apply attention mask
                sum_hidden = torch.sum(hidden * attention_mask_expanded, dim=1)  # Sum over all tokens
                count_non_padding = attention_mask.sum(dim=1, keepdim=True)  # Count of non-padding tokens
                vector = sum_hidden / count_non_padding  # Average embedding
                vector = vector.detach().cpu().flatten()  # Flatten the vector
                database = base_database

            distances = cosine_similarity([vector], database)
            distances = list(float(x) for x in distances.flatten())
            result = list(zip(df.text, distances))
            result_sorted = sorted(result, key=lambda x: x[1], reverse=True)

            print(f"Query: {text}")
 
            for k in range(3):
                print(f"\nResult {k} (L2 = {round(result_sorted[k][1], 5)}):")
                print(f"Text: {result_sorted[k][0]}")

