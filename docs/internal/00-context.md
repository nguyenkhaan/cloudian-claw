
# 1. Context là gì 
- context là một pakcage STD của Go. Dùng để mang trạng thái của một operation. 

Các trạng thái thông tin cơ bản mà Context có thể quản lý: 
    + Cancellation: operation đã bị hủy chưa 
    + Timeout: operation chạy tối đa trong bao lâu 
    + Deadline: Operation phải kết thúc trước thời điểm nào 
    + Value: operation có những dữ liệu nào 

Interface cơ bản: 
```go 
type Context interface {
    Deadline() (deadline time.Tine, ok book) 
    Done() <-chan struct{} 
    Err() error 
    Value(key any) any 
}
```

## Context có quan hệ Parent -> Child 
Context thường được tạo ra dựa trên một context khác 

**Ví dụ**: 
```go
ctx := context.Background() 
child, cancel := context.WithCancel(ctx) 
```

```go
ctx1 := context.Background()

ctx2, cancel := context.WithCancel(ctx1)
ctx3, cancel := context.WithTimeout(ctx2, 5*time.Second)
```

## context.Background() 
- Đây là context gốc, giống như một context rỗng. Nó thường đóng vai trò như 1 điểm bắt đầu, để tạo ra các context khác. 

Chú ý: context.Background() vẫn chứa 4 API như 1 context bình thường 

## context.TODO() 
- context TODO cũng tạo ra một context không có thông tin đặc biệt. Tuy nhiên, khác biệt chính nằm ở ý nghĩa: 
    + background(): Cái này là gốc 
    + TODO(): Tôi chưa xác định được context thích hợp để dùng -> thường dùng trong quá trình phát triển


# 2. Các loại context 
## context.WithCancel() 
API này tạo ra một context có khả năng cancel chủ động. 

Nó sẽ tự động cancel(). Trong code, chúng ta có 1 hàm để bắt được sự kiện này, và ngay lập tức chấm dứt công việc. Vì vậy, API này được dùng với các công việc, chỉ được phép thực hiện trong 1 khoảng tgian nhất định. 


**Cú pháp**: 
```go 
ctx, cancel := context.WithCancel(parent)
```
Hàm cancel được bạn giữ, bạn sẽ gọi cancel() khi muốn kết thúc context này. 

**Ví dụ**: 
```go 
func main() {
    ctx, cancel := context.WithCancel(ctx.Background())
    //Mot tgian nao do 
    cancel() 
}
```

## context.WithTimeOut() 
Tạo một Context có giới hạn thời gian 

**Cú pháp** 
```go
ctx, cancel := context.WithTimeOut(parent , time)  
```

Sau tgian, context sẽ tự động chuyển qua trạng thái hết hạn. Bạn cũng có thể dùng cancel() để chủ động kết thúc trước thời hạn. 

## context.WithDeadline() 
- Cũng giới hạn thời gian, nhưng thay vì truyền khoảng thời gian, bạn truyền một thời điểm cụ thể. 

```go
ctx, cancel := context.WithDeadline(context.Background(), deadline,) 
```
## context.WithValue() 
Tạo một context mới, chứa một cặp key - value. Dùng để truyền giá trị giữa các hàm. 

**Cú pháp** 
```go
ctx := context.WithValue(context.Background() , key, value,) 
```
Để lấy ra được giá trị, chúng ta có thể dùng ctx.Value(key) 
