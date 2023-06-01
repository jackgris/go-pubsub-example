const Message: React.FC<{ index:number,  username: string, message: string }> = ({ index, username, message }) => {
       // return <div key={index}> {username}: {message} </div>
return <div key={index} className="flex items-center bg-blue-500 text-white rounded-lg p-2 mb-2">
  <span className="font-semibold">{username}:</span>
  <span className="ml-2">{message}</span>
</div>

}

export default Message;
