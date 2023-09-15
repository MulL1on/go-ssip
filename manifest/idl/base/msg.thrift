namespace go msg

struct Msg {
    1: i8 type
    2: i64 from_user
    3: i64 to_user
    4: i64 to_group
    5: string content
}

struct  GroupMsg{
    1: i64 from_user
    2: i64 group
    3: string content
}