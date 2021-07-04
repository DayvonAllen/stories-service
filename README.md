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
---