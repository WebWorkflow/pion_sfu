const pionSFU = {
    socket: WebSocket,
    pc: RTCPeerConnection,
    room_id: String,
    self_id: String
}



pionSFU.connect = () => {
    console.log("Connected");

    return new WebSocket('ws://localhost:8080')
}
pionSFU.makeId = (length) => {
    let result = '';
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    const charactersLength = characters.length;
    let counter = 0;
    while (counter < length) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
        counter += 1;
    }
    return result;
}
pionSFU.sendOffer = async (socket, pc, self_id, room_id) => {
    const offer = await pc.createOffer();
    await pc.setLocalDescription(offer);
    const message= {
        self_id: self_id,
        room_id: room_id,
        offer: offer
    }
    console.log('Offer created: ', message);

    socket.send(JSON.stringify({event: 'offer', data: message}));
}


pionSFU.joinRoom = async () => {
    const message= {
        room_id: pionSFU.room_id,
        self_id: pionSFU.self_id
    };
    pionSFU.socket.send(JSON.stringify({ event: 'joinRoom', data: message }));
}

pionSFU.leaveRoom = async () => {
    const message= {
        room_id: pionSFU.room_id,
        self_id: pionSFU.self_id
    }
    pionSFU.socket.send(JSON.stringify({ event: 'leaveRoom', data: message }));
}


pionSFU.setup = () => {
    pionSFU.socket.onmessage = async (event) => {
        const data = JSON.parse(event.data.toString());
        switch (data.event) {
            case 'answer': {
                const value = JSON.parse(data.data);
                pionSFU.pc.setRemoteDescription(value);
                pionSFU.pc.oniceconnectionstatechange = (event) => {
                    console.log(event);
                    console.log(pionSFU.pc.iceConnectionState);
                };
                pionSFU.pc.onicegatheringstatechange = (event) => {
                    console.log(event);
                    console.log(pionSFU.pc.iceConnectionState);
                }
                break;
            }
            case 'candidate': {
                const value = JSON.parse(data.data);
                await pionSFU.pc.addIceCandidate(new RTCIceCandidate(value));
                break;
            }
            case 'offer': {
                const value = JSON.parse(data.data);
                console.log("Offer ", value);
                await pionSFU.pc.setRemoteDescription(value);
                await pionSFU.pc.createAnswer().then(
                    async (answer) => {
                        await pionSFU.pc.setLocalDescription(answer);
                        const message = {
                            self_id: value.self_id,
                            room_id: value.room_id,
                            answer: answer
                        }
                        pionSFU.socket.send(JSON.stringify({ event: 'answer', data: message }));
                    }
                );
                break;
            }
        }
    };

    pionSFU.pc.onicecandidate = (event) => {
        event.candidate ? pionSFU.socket.send(JSON.stringify({ event: 'ice-candidate', data: { room_id: pionSFU.room_id, self_id: pionSFU.self_id, candidate: event.candidate } })) : null;
    };
};




pionSFU.init = async (room_id) => {
    pionSFU.room_id = room_id
    pionSFU.self_id = pionSFU.makeId(10)
    pionSFU.socket = pionSFU.connect()
    await pionSFU.joinRoom()
    pionSFU.pc = new RTCPeerConnection()
    pionSFU.setup()
    await pionSFU.sendOffer()
}