namespace go msg

struct Msg {
    1: i8 type
    3: i64 from_user
    4: i64 to_user
    5: i64 to_group
    6: string text
}

struct  GroupMsg{
    1: i64 from_user
    2: i64 group
    3: string content
}