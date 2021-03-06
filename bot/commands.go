package bot

import (
	"fmt"
	"strings"

	"github.com/thehowl/setabot/stops"
	"gopkg.in/telegram-bot-api.v4"
)

const startMessage = `Ciao! Sono un bot per vedere quanto manca all'arrivo del tuo autobus, per le provincie di Modena, Reggio Emilia e Piacenza.
Sono ancora in "beta" testing circa, se vuoi chiedere qualcosa al mio sviluppatore scrivigli qui: @dahhowl.
Per chiedermi gli autobus che devono arrivare ad una fermata, usa il comando /qm, ad esempio:
/qm Modena Autostazione
/qm Marzaglia Vecchia
/qm Ca' Bianca

Il bot può essere un po' stupido a volte a capire a che fermata ti stai riferendo. Prima o poi lo metterò a posto.
Sorgente disponibile qui: https://github.com/thehowl/setabot`

// welcome message
func (b *Bot) start(u tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
	msg.Text = startMessage
	if b.getCity(u.Message.From.ID) == "" {
		msg.Text += "\n\nAd ogni modo, mi servirebbe sapere di che città sei."
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("Sono di Modena"),
			tgbotapi.NewKeyboardButton("Sono di Reggio Emilia"),
			tgbotapi.NewKeyboardButton("Sono di Piacenza"),
		})
	}
	b.send(msg)
}

// register that the user is from a certain city.
func (b *Bot) imFrom(u tgbotapi.Update) {
	abbr := cities[strings.ToLower(u.Message.Text)]
	if abbr == "" {
		return
	}
	b.setCity(u.Message.From.ID, abbr)
	b.send(tgbotapi.NewMessage(u.Message.Chat.ID, "Ok! Ho registrato che vieni da "+u.Message.Text))
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Se devi cambiare la tua provincia, ti basterà dire 'Sono di x', e se scrivi la città bene cambierò la tua città.")
	msg.ReplyMarkup = tgbotapi.NewHideKeyboard(false)
	b.send(msg)
}

// Quanto Manca?™
func (b *Bot) qm(u tgbotapi.Update) {
	// Get the user's city, and forbit getting quanto manca if there's no city
	// set
	city := b.getCity(u.Message.From.ID)
	if city == "" {
		b.send(tgbotapi.NewMessage(u.Message.Chat.ID, "Mi serve prima sapere di che città sei!"))
		return
	}

	// iterate over the stops of our city to find one that looks like the name.
	cStops := stops.CityStops[city]
	stopName := u.Message.Text
	var chosen *stops.Stop
	for _, stop := range cStops {
		if strings.ToLower(stopName) == strings.ToLower(stop.Name) {
			chosen = &stop
			break
		}
	}

	// no stop found :(
	if chosen == nil {
		b.send(tgbotapi.NewMessage(u.Message.Chat.ID, "Mi dispiace, purtroppo non conosco questa fermata!"))
		return
	}

	// send chat action since this will take long and we don't want to look dead
	b.send(tgbotapi.NewChatAction(u.Message.Chat.ID, "typing"))

	// get arrivals and check they're there
	arrivals, err := b.AS.GetArrivals(city, chosen.Identifier, chosen.Name)
	if err != nil {
		fmt.Println(err)
		b.send(tgbotapi.NewMessage(u.Message.Chat.ID, "C'è stato un errore! Non è stato possibile vedere gli arrivi di questa fermata."))
		return
	}

	if len(arrivals) == 0 {
		b.send(tgbotapi.NewMessage(u.Message.Chat.ID, "Hmm... Non sembra esserci nessun autobus a quest'ora!"))
		return
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
	msg.Text = fmt.Sprintf("Ho trovato %d risultati per %s:\n\n", len(arrivals), chosen.Name)
	for _, arrival := range arrivals {
		msg.Text += fmt.Sprintf(
			"* <b>%s</b> in direzione %s, arriva tra <b>%s minuti</b> (alle %s)\n",
			arrival.Line, arrival.Destination, arrival.ToArrival, arrival.RealTime.Format("15:04"),
		)
	}

	msg.ParseMode = "HTML"

	b.send(msg)
}

// search through stops for something with that name.
func (b *Bot) search(u tgbotapi.Update) {
	city := b.getCity(u.Message.From.ID)
	if city == "" {
		b.send(tgbotapi.NewMessage(u.Message.Chat.ID, "Mi serve prima sapere di che città sei!"))
		return
	}

	// get the stops for the city
	sts := stops.CityStops[city]

	// n is the number of stops currently matched, and stopsText is the message
	// to return.
	var (
		n         int
		stopsText string
	)

	// ToLower because we want the search to be case insensitive
	u.Message.Text = strings.ToLower(u.Message.Text)

	for _, st := range sts {
		if strings.Contains(strings.ToLower(st.Name), u.Message.Text) {
			// if the matched result is the 151st, then we just skip and say
			// there are more than
			if n == 150 {
				stopsText += "* E altri...\n"
				break
			}
			stopsText += "* " + st.Name + "\n"
			n++
		}
	}

	b.send(tgbotapi.NewMessage(u.Message.Chat.ID, fmt.Sprintf("Ottenuti %d risultati.\n%s", n, stopsText)))
}
