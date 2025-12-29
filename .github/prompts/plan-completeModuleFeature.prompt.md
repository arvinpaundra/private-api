# Plan: Complete Module Feature with DDD Infrastructure Layer

Implement the complete module feature following DDD principles, including infrastructure repositories for module operations and context mappers for grade/subject validation. Fix service error handling and create the presentation layer (handler/router) following established patterns.

## Steps

1. **Create infrastructure layer** in `infrastructure/module/` with three repositories: `ModuleWriterRepository` implementing `Save()` with insert/update/remove branching (following `infrastructure/grade/grade_writer_repository.go` pattern), `GradeMapperRepository` implementing `IsGradeExist()` with EXISTS query scoped by user and soft-delete filtering, and `SubjectMapperRepository` implementing `IsSubjectExist()` with the same pattern.

2. **Add domain response DTOs** in `domain/module/response/` creating `module.go` with `Module` struct containing fields: `ID`, `UserID`, `SubjectID`, `GradeID`, `Title`, `Slug`, `Description`, `Type`, `IsPublished`, `Questions` (following `domain/grade/response/grade.go` structure), and optionally `question.go`/`question_choice.go` for nested responses.

3. **Fix service error handling** in `domain/module/service/create_module.go` by changing the two `return nil` statements (lines 41, 48) to return `constant.ErrSubjectNotFound` and `constant.ErrGradeNotFound` respectively, ensuring proper error propagation to the handler.

4. **Create presentation layer** with `application/rest/handler/module.go` implementing `ModuleHandler` struct and `CreateModule` method (following `handler/grade.go` pattern with validation, auth storage, repository instantiation, and error mapping), and `application/rest/router/module/router.go` with `ModuleRouter` struct and `Private()` method registering POST `/modules` endpoint with authentication middleware.

5. **Wire up module router** in `application/rest/router/router.go` by instantiating `ModuleRouter` and calling its `Private()` method within the v1 API group (following the pattern used for grade/subject routers).

## Further Considerations

1. **Context Mapping Strategy**: The mapper repositories implement the **Conformist** pattern (grade/subject are upstream contexts, module conforms to their interfaces). Should we add anti-corruption layer validation or keep the simple EXISTS checks? Consider adding caching for frequently-checked subject/grade existence.

Answer:
Just keep simple EXISTS checks for now to minimize complexity. Caching can be considered later if performance issues arise.

2. **Error Handling for Mappers**: Should mapper repositories return domain-specific errors (`constant.ErrSubjectNotFound`) or generic GORM errors? Current grade/subject infrastructure returns raw GORM errors, but service layer could benefit from domain errors for clearer error mapping in handlers.

Answer:
Mapper repositories should return generic errors; the service layer can interpret these and return domain-specific errors as needed.

3. **Response DTO Scope**: Should `ModuleResponse` include nested `Subject` and `Grade` response objects for richer API responses, or keep IDs only for simplicity? Also, will CRUD operations beyond Create be needed (Update, Delete, Get, List)?

Answer:
Keep it simple with IDs only for now. Additional CRUD operations can be added later as needed.
