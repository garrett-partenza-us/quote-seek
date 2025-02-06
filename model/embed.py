import torch
import pandas as pd
from encoder import Encoder
from tqdm import tqdm
from transformers import BertTokenizer, BertModel
from transformers import BertForSequenceClassification, BertTokenizer

model_name = "bert-base-uncased"
meditations = "../data/meditations.txt"

fine_tuned_model = BertForSequenceClassification.from_pretrained("./fine_tuned_bert")
base_model = BertModel.from_pretrained("bert-base-uncased")

tokenizer = BertTokenizer.from_pretrained("./fine_tuned_bert")

device = torch.device("mps" if torch.backends.mps.is_available() else "cpu")
fine_tuned_model.to(device)
base_model.to(device)

encoder = Encoder(model_name)
encodings, sentence_ids = encoder(meditations)

batch_size = 16
embeddings = []

print(len(sentence_ids))

for batch_num in tqdm(range(0, len(sentence_ids), batch_size)):

    input_ids = encodings['input_ids'][batch_num:min(batch_num + batch_size, len(sentence_ids))]
    attention_mask = encodings['attention_mask'][batch_num:min(batch_num + batch_size, len(sentence_ids))]
    current_sentence_ids = sentence_ids[batch_num:min(batch_num + batch_size, len(sentence_ids))]
    texts = [encoder.lookup_index(x) for x in current_sentence_ids]

    with torch.no_grad():
        fine_tuned_outputs = fine_tuned_model.bert(input_ids=input_ids.to(device),
                                                      attention_mask=attention_mask.to(device))
        fine_tuned_hidden_states = fine_tuned_outputs.last_hidden_state
        fine_tuned_vectors = fine_tuned_hidden_states[:, 0, :]

        base_outputs = base_model(input_ids=input_ids.to(device), attention_mask=attention_mask.to(device))
        base_hidden_states = base_outputs.last_hidden_state
        base_vectors = base_hidden_states[:, 0, :]

    for embedding_id, item in enumerate(current_sentence_ids):
        embeddings.append({
            "text": texts[embedding_id],
            "fine_tuned_embedding": fine_tuned_vectors[embedding_id].detach().cpu().numpy(),
            "base_embedding": base_vectors[embedding_id].detach().cpu().numpy(),
        })


df = pd.DataFrame(embeddings)
df.to_pickle("./vectors.pkl")

print("Vectors saved successfully.")
