import 'dotenv/config'
import 'module-alias/register'
import { Mailer } from '../application/mail/mailer'
import { MessageQueue } from '../application/message-queue/message-queue'
import { BullMailQueue } from '../framework/mail-queue/bull-mail-queue'
import { NodemailerMailProvider } from '../framework/mail/nodemailer-mail-provider'
;(async () => {
  try {
    const mailer = new Mailer(new NodemailerMailProvider())
    const mailQueue = new BullMailQueue(mailer)

    // start consuming message queue
    const queue = new MessageQueue(mailQueue)
    await queue.consume()

    // start processing email submissions
    mailQueue.handleFailedJobs()
    await mailQueue.process()
  } catch (error) {
    console.error(error)
    process.exit(1)
  }
})()
