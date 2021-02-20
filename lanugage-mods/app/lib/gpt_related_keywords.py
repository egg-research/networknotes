# Inspired by https://github.com/shreyashankar/gpt3-sandbox/tree/master 
from .gpt import GPT, Example

def getRelatedKeywordGPT():
    # Prime GPT object to produce related keywords
    gpt = GPT(engine="curie",
            temperature=0.5,
            max_tokens=100)

    gpt.add_example(Example('cerebellum, motor control, sensorimotor calibration, purkinje cell',
                            'motor learning, thin neuronal layers, synaptic plasticity, parallel connections, granule cell'))

    gpt.add_example(Example("amygdala, hemispheric emotions, fear, independent memory systems, declarative memory",
                            'ptsd, memory modulation, lesions, emotion learning, synapses, working memory'))

    gpt.add_example(Example("broca's area, removal, speech impairment, broken grammar, aphasia",
                            "expressive aphasia, speech generation, phenomic linking, language, wernicke's area"))

    gpt.add_example(Example('create process, control, operating system, information transfer, software extension of hardware structure',
                            'system nucleus, primatives, synchronization, multiprogramming system, remove process'))

    gpt.add_example(Example('semaphore, avoid race conditions, arbitrary resource count, binary locks, queue', 
                            'monitor, Edsger Dijkstra, starvation, first in first out, decrement'))
                            
    gpt.add_example(Example('distributed systems, parallel computing, distributed computing, local memory, communication links', 
                            'networked, HTTP, lack global clock, achieve common goal, massively multiplayer online games',))

    return gpt