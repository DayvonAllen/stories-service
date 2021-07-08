## Routes(All routes are protected except for featured)
- Get Featured stories(only returns 10 stories):
  - `GET:http://localhost:8081/stories/featured`
- Get All stories:
    - `GET:http://localhost:8081/stories?page=1`(Gives 10 stories at a time)`
- Get All stories(returns new stories first):
  - `GET:http://localhost:8081/stories?page=1&new=true`(Gives 10 stories at a time)`
- Get story by ID:
  - `GET:http://localhost:8081/stories/<story ID>`
- Create Story:
    - `POST:http://localhost:8081/stories`
    - JSON(must include at least one tag):
```
{
    "title": "My unknown title",
    "content": "my unknown content",
    "tags": [
        {
            "value": "CreepyPasta"
        },
          {
            "value": "TrueScaryStory"
        },
          {
            "value": "CampFire"
        },
          {
            "value": "Paranormal"
        },
          {
            "value": "GhostStory"
        },
        {
            "value": "Other"
        }
    ]
}
```
- Delete Story:
    - `DELETE:http://localhost:8081/stories/<story ID>`
- Update Story:
    - `PUT:http://localhost:8081/stories/<Story ID>`
    - JSON:
```
{
    "content": "updated content",
    "title": "updated title",
    "tags": [
          {
            "value": "Paranormal"
        },
          {
            "value": "GhostStory"
        },
        {
            "value": "Other"
        }
    ]
}
```    
- Like Story:
    - `PUT:http://localhost:8081/stories/like/<Story ID>`
- Dislike Story:
    - `PUT:http://localhost:8081/stories/dislike/<Story ID>`
- Like Comment:
    - `http://localhost:8081/comment/like/<Comment ID>`
- Dislike Comment:
    - `http://localhost:8081/comment/dislike/<Comment ID>`
- Like Reply:
  - `http://localhost:8081/reply/like/<Reply ID>`
- Dislike Reply:
  - `http://localhost:8081/reply/dislike/<Reply ID>`  
- Comment On A Story:
  -`POST:http://localhost:8081/comment/<Story ID>`
  - JSON:
```
{
    "content": "Nice story"
}
```  
- Update comment:
  -`PUT:http://localhost:8081/comment/<Comment ID>`
  - JSON:
```
{
    "content": "Nice story"
}
```  
- Update Reply:
  -`PUT:http://localhost:8081/reply/<Reply ID>`
  - JSON:
```
{
    "content": "Nice story"
}
```  
- Delete Comment:
  -`DELETE:http://localhost:8081/comment/<Comment ID>`
- Delete Reply:
  -`DELETE:http://localhost:8081/reply/<Reply ID>`  
- Flag Story:
  - `PUT:http://localhost:8081/stories/flag/<Story ID>`
- Flag Comment:
  - `PUT:http://localhost:8081/comment/flag/<Comment ID>`
- Flag Reply:
  - `PUT:http://localhost:8081/reply/flag/<Flag ID>`  
- Reply To Comment:
  - `POST: http://localhost:8081/reply/<Comment ID>`
  - JSON:
```
{
    "content": "reply"
}
```  
- Get all read later items by username
  - `GET:http://localhost:8081/read`
- Add Story to read later list:
  - `POST:http://localhost:8081/read/<Story ID>`
- Delete Story from read later list:
  - `DELETE:http://localhost:8081/read/<Story ID>`
---