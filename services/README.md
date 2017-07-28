Api contract
--

Listen for HTTP POST on 80 for `/execute`. 

POST body has structure 
```
{
  'name': 'foo',
  'args': 'bar fizz',
}
```


Must respond with response object with structure
```
{
  'command': {
    'name': 'foo',
    'args': 'bar fizz'
  },
  'type': 'success|failure',
  'answers': [
    'buzz',
    'bizz',
    'bap',
  ],
}
```
