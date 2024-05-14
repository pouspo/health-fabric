## Dummy Policies
```json
{
  "id": "group1",
  "policy_map": {
    "alpha": {
      "read": [
        "glucose",
        "bmi",
        "age"
      ],
      "write": [
        "glucose",
        "bmi"
      ]
    },
    "beta": {
      "read": [
        "glucose",
        "bmi"
      ],
      "write": [
        "glucose",
        "bmi"
      ]
    },
    "gama": {
      "read": [
        "glucose",
        "bmi"
      ],
      "write": []
    }
  }
}
```
```json
{
  "id": "group2",
  "policy_map": {
    "alpha": {
      "read": [
        "skin_thickness"
      ],
      "write": [
        "skin_thickness"
      ]
    },
    "beta": {
      "read": [
        "skin_thickness"
      ],
      "write": [
        "skin_thickness"
      ]
    }
  }
}

```