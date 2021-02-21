from google.cloud import language_v1
import re

def top_keywords(text_content, num_keywords=10):
    """
    Analyzing Entities in a String

    Args:
      text_content The text content to analyze
    """

    client = language_v1.LanguageServiceClient()
    type_ = language_v1.Document.Type.PLAIN_TEXT

    language = "en"
    document = {"content": text_content, "type_": type_, "language": language}

    # Available values: NONE, UTF8, UTF16, UTF32
    encoding_type = language_v1.EncodingType.UTF8
    
    # Get response
    response = client.analyze_entities(request = {'document': document, 'encoding_type': encoding_type})
    keywords = set()
    # Process each response entity to create a list of 10 unique keywords
    for entity in response.entities:
        
        new_word = entity.name.lower().strip()
        word_already_in_set = False
        for kwd in keywords: 
            if re.search('^' + kwd + '(s)?(ing)?(ed)?$', new_word) is not None: 
                word_already_in_set = True
        if not word_already_in_set: 
            keywords.add(new_word)
        
        if len(keywords) > num_keywords: 
            break
    return list(keywords)