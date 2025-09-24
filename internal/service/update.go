package service

func (s *Service) UpdateBookingStatus(id int, newStatus string) error {
	return s.storage.UpdateBookingStatus(id, newStatus)
}
