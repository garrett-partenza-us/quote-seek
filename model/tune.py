import os
import torch
from transformers import BertTokenizer, BertForMaskedLM
from torch.utils.data import DataLoader, TensorDataset
from torch.optim import AdamW
from nltk.tokenize import sent_tokenize
from encoder import Encoder

model_name = 'bert-base-uncased'
meditations = '../data/meditations.txt'

# Check if the M1 GPU (MPS device) is available
device = torch.device("mps" if torch.backends.mps.is_available() else "cpu")
print(f"Using device: {device}")

# Initialize model
model = BertForMaskedLM.from_pretrained(model_name)
model.to(device)

# Freeze all layers except the last few
for name, param in model.named_parameters():
    if any(x in name for x in ['encoder.layer.9', 'encoder.layer.10', 'encoder.layer.11']):
        param.requires_grad = True
    else:
        param.requires_grad = False

for name, param in model.named_parameters():
    print(f"Layer: {name}, Trainable: {param.requires_grad}")

# Initialize Encoder
encoder = Encoder(model_name)

# Train setup
num_epochs = 50
batch_size = 128
learning_rate = 5e-5

optimizer = AdamW(filter(lambda p: p.requires_grad, model.parameters()), lr=learning_rate)

# Train loop
model.train()

for epoch in range(num_epochs):
    # Re-encode and remask data at the beginning of each epoch
    encodings, sentence_ids = encoder(meditations)
    inputs, labels = encoder.mask_tokens(encodings['input_ids'])

    # Prepare DataLoader for the epoch with new masked data
    train_dataset = TensorDataset(encodings['input_ids'], encodings['attention_mask'], labels)
    train_dataloader = DataLoader(train_dataset, batch_size=batch_size, shuffle=True)

    # Iterate over the batches for the current epoch
    for step, batch in enumerate(train_dataloader):
        batch_input, batch_attention_mask, batch_labels = batch
        optimizer.zero_grad()

        outputs = model(
            batch_input.to(device),
            attention_mask=batch_attention_mask.to(device),
            labels=batch_labels.to(device)
        )

        loss = outputs.loss
        loss.backward()
        optimizer.step()

        # Print the loss at each step
        print(f"Epoch {epoch+1}, Batch {step+1}, Loss: {loss.item()}")

# Save the fine-tuned model and tokenizer
model.save_pretrained('./fine_tuned_bert')
encoder.tokenizer.save_pretrained('./fine_tuned_bert')

