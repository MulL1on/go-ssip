namespace go msg

struct Msg {
    1: i64 from_user
    2: i64 to_user
    3: string content
}

struct  GroupMsg{
    1: i64 from_user
    2: i64 group
    3: string content
}