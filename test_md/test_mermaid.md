# Mermaid 圖表測試

## 流程圖 (Flowchart)

```mermaid
flowchart LR
    A[Start] --> B{Is it?}
    B -->|Yes| C[OK]
    B -->|No| D[End]
    C --> D
```

## 序列圖 (Sequence Diagram)

```mermaid
sequenceDiagram
    participant Alice
    participant Bob
    Alice->>Bob: Hello Bob!
    Bob->>Alice: Hi Alice!
```

## 圓餅圖 (Pie Chart)

```mermaid
pie title Languages
    "Go" : 40
    "Swift" : 30
    "Python" : 20
    "Other" : 10
```

## 甘特圖 (Gantt Chart)

```mermaid
gantt
    title Project Timeline
    dateFormat YYYY-MM-DD
    section Design
    Design Phase :a1, 2026-01-01, 30d
    section Development
    Development Phase :a2, after a1, 60d
    section Testing
    Testing Phase :a3, after a2, 20d
```

## 類圖 (Class Diagram)

```mermaid
classDiagram
    class Animal {
        +String name
        +makeSound()
    }
    class Dog {
        +bark()
    }
    class Cat {
        +meow()
    }
    Animal <|-- Dog
    Animal <|-- Cat
```

## 狀態圖 (State Diagram)

```mermaid
stateDiagram-v2
    [*] --> Idle
    Idle --> Running : start
    Running --> Paused : pause
    Paused --> Running : resume
    Running --> [*] : stop
```

## 使用者旅程圖 (User Journey)

```mermaid
journey
    title Shopping Journey
    section Browse
      Search for item :5: User
      View details :4: User
    section Purchase
      Add to cart :3: User
      Checkout :2: User
```

## 需求圖 (Requirement Diagram)

```mermaid
requirementDiagram
    requirement test_req {
        id: 1
        text: the test text
        risk: low
        verifymethod: analysis
    }
    element test_entity {
        type: simulation
    }
    test_entity - satisfies -> test_req
```

## XY 圖表 (XY Chart)

```mermaid
xychart-beta
    title "Sales Revenue 2024"
    x-axis [Jan, Feb, Mar, Apr, May, Jun]
    y-axis "Revenue (USD)" 4000 --> 12000
    bar [5000, 6000, 7500, 8200, 9800, 10500]
    line [5000, 6000, 7500, 8200, 9800, 10500]
```

## XY 圖表 (水平)

```mermaid
xychart-beta horizontal
    title "Browser Market Share"
    x-axis [Chrome, Firefox, Safari, Edge]
    y-axis "Percentage" 0 --> 100
    bar [65, 20, 10, 5]
```