package task

type TaskServiceImpl struct {
	repository TaskRepository
}

func NewTaskService(repo TaskRepository) (TaskService) {
	return&TaskServiceImpl{
		repository: repo,
	}
}

func (s *TaskServiceImpl) HandleAdd(newTask NewTaskSchema) (error) {

}

func (s *TaskServiceImpl) HandleList(completed bool, all bool) ([]Task, error) {

}

func (s *TaskServiceImpl) HandleDelete(id int) (error) {

}

func (s *TaskServiceImpl) HandleDone(id int) (error) {

}