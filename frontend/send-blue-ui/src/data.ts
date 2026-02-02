// filepath: c:\Users\lastb\Code\send_blue\frontend\send-blue-ui\src\types.d.ts

export interface Client {
    uid: string;
    name: string;
    messagesSent: number;
    leads: ClientLead[];
    messageQueue: MessageQueue[];
    allMessagesSent: MessageQueue[];
}

export interface ClientLead {
    leadNumber: string;
    clientUid: string;
    messagesReceived: number;
    lastContacted: Date | null;

}

export interface MessageQueue {
    msgUid: string;
    messageBody: string;
    fromClientId: string;
    toClientLead: string;
    scheduledSendTime: Date | null;
    timeSent: Date | null;
    status: string;
}
// Dashboard message event history interface
export interface MessageEventHistory {
    msgUid: string;
    time_stamp: Date;
    prev_status: string | null;
    curr_status: string;
}

/////////////////API CALLS//////////////////////

export async function fetchClients(): Promise<Client[]> {
    console.log('fetchClients called');
    const clients: Client[] = []
    try {
        const base = import.meta.env.VITE_API_BASE || ''    
        const res = await fetch(`${base}/clients`)
        if (!res.ok) throw new Error('Network')
        const data = await res.json()  
        clients.push(...(data.clients || []))
      } catch (e) {
        console.error('fetchClients error', e)
      } finally {
        return clients
      }
}

export async function fetchClientLeads(clientID: string): Promise<ClientLead[]> {
    const leads: ClientLead[] = []
     try {
        const base = import.meta.env.VITE_API_BASE || ''
        const res = await fetch(`${base}/clients/leads?client_id=${clientID}`)
        if (!res.ok) throw new Error('Network')
        const data = await res.json()
        leads.push(...(data.leads || []))
      } catch (e) {
        console.error('fetchClientLeads error', e)
      } finally {
        return leads
      }
}
export async function fetchMessageQueue(clientID: string): Promise<MessageQueue[]> {
    const queue: MessageQueue[] = []
        try {
        const base = import.meta.env.VITE_API_BASE || '' 
        const res = await fetch(`${base}/clients/scheduled?client_id=${clientID}`)
        if (!res.ok) throw new Error('Network')
        const data = await res.json()
        console.log('Fetched scheduled messages:', data)
        queue.push(...(data.messages || []))
      } catch (e) {
        console.error('fetchMessageQueue error', e)
      } finally {
        return queue
      }
    }
export async function fetchClientData(clientID: string): Promise<Client> {
    const client: Client = { uid: '', name: '', messagesSent: 0, messageQueue: [], allMessagesSent: [], leads: []  }
    try {
        const base = import.meta.env.VITE_API_BASE || ''
        const res = await fetch(`${base}/client/data?client_id=${clientID}`)
        if (!res.ok) throw new Error('Network')
        const data = await res.json()
        client.uid = data.client.uid
        client.name = data.client.name
        client.messagesSent = data.client.messagesSent
        client.leads = data.client.leads || []
        client.messageQueue = data.client.messageQueue || []
        client.allMessagesSent = data.client.allMessagesSent || []
      } catch (e) {
        console.error('fetchClientData error', e)
      } finally {
        return client
      }
}

