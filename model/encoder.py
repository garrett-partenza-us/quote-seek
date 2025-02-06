import torch
from transformers import BertTokenizer
from nltk.tokenize import sent_tokenize

class Encoder():

    def __init__(self, model: str, window_size: int = 64, overlap: int = 32):
        """
        Initialize the encoder with a specific window size and overlap.
        Args:
            model (str): The transformer model name to load the tokenizer.
            window_size (int): Maximum number of tokens per window.
            overlap (int): Number of tokens to overlap between consecutive windows.
        """
        self.tokenizer = BertTokenizer.from_pretrained(model)
        self.max_token_length = self.tokenizer.model_max_length
        self.sentences = None
        self.window_size = window_size
        self.overlap = overlap

    def __call__(self, path: str):
        self.sentences = self.get_sentences(path)
        sentence_indices = []
        tokenized_data = []

        for idx, sentence in enumerate(self.sentences):
            chunks = self.chunk_sentence(sentence)
            for chunk in chunks:
                tokenized_data.append(chunk)
                sentence_indices.append(idx)

        return self.pad_tokens(tokenized_data), sentence_indices

    def encode_query(self, query: str):
        self.sentences = [query]
        sentence_indices = []
        tokenized_data = []

        chunks = self.chunk_sentence(query)
        for chunk in chunks:
            tokenized_data.append(chunk)
            sentence_indices.append(0)

        return self.pad_tokens(tokenized_data), sentence_indices

    def get_sentences(self, path):
        sentences = []
        with open(path) as file:
            for entry in file.readlines():
                entry = entry.strip()
                if entry:
                    for sentence in sent_tokenize(entry):
                        sentence = sentence.strip()
                        if sentence:
                            sentences.append(sentence.strip())
        return sentences

    def chunk_sentence(self, sentence):
        tokens = self.tokenizer.encode(sentence, add_special_tokens=False)
        # If the sentence is too long, split it into smaller windows using a sliding window approach
        if len(tokens) > self.window_size:
            chunks = []
            for i in range(0, len(tokens), self.window_size - self.overlap):
                chunk = tokens[i:i + self.window_size]
                chunks.append(chunk)
            return chunks
        else:
            return [tokens]

    def pad_tokens(self, tokenized_data):
        # Create a dictionary to hold the tokenized data
        tokenized_dict = {
            'input_ids': tokenized_data,
            'attention_mask': [[1] * len(tokens) for tokens in tokenized_data]  # Create attention mask
        }

        # Pad the tokenized data and return the result
        padded_encodings = self.tokenizer.pad(
            tokenized_dict,
            padding=True,
            return_tensors="pt",
            return_attention_mask=True
        )
        return padded_encodings

    def lookup_index(self, idx: int):
        assert self.sentences, "Encoder not called"
        assert 0 <= idx <= len(self.sentences), "Index out of range"
        return self.sentences[idx]

    def mask_tokens(self, inputs, mlm_probability=0.15):
        labels = inputs.clone()
        mask = torch.bernoulli(torch.full(labels.shape, mlm_probability)).bool()
        inputs[mask] = self.tokenizer.mask_token_id
        labels[~mask] = -100
        return inputs, labels

