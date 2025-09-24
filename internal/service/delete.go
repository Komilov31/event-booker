package service

func (s *Service) DeleteEvent(id int) error {
	return s.storage.DeleteEvent(id)
}

func (s *Service) DeleteBooking(id int, placeCount int) error {
	return s.storage.DeleteBooking(id, placeCount)
}
