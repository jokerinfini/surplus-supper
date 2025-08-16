# 🤝 Contributing to Surplus Supper

Thank you for your interest in contributing to Surplus Supper! This document provides guidelines and information for contributors.

## 🎯 About Surplus Supper

Surplus Supper is a marketplace platform that connects restaurants with customers to reduce food waste by selling surplus food at discounted prices. Our mission is to create a sustainable food ecosystem while helping restaurants reduce waste and customers save money.

## 🚀 Quick Start

### Prerequisites
- **Go** 1.21+ (for backend)
- **Node.js** 18+ (for frontend)
- **Docker** (for database)
- **Git**

### Setup Development Environment

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/surplus-supper.git
   cd surplus-supper
   ```

2. **Start Database**
   ```bash
   docker-compose up -d postgres
   ```

3. **Setup Backend**
   ```bash
   cd backend
   go mod download
   go run main.go
   ```

4. **Setup Frontend**
   ```bash
   cd frontend-next
   npm install
   npm run dev
   ```

5. **Verify Setup**
   - Backend: http://localhost:8080/health
   - Frontend: http://localhost:3000
   - Database: PostgreSQL on port 5433

## 📋 How to Contribute

### 1. Find an Issue
- Check the [Issues](https://github.com/YOUR_USERNAME/surplus-supper/issues) page
- Look for issues labeled `good first issue` for beginners
- Comment on an issue you'd like to work on

### 2. Create a Branch
```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 3. Make Changes
- Follow the coding standards below
- Write tests for new features
- Update documentation as needed

### 4. Test Your Changes
```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend-next
npm run test
```

### 5. Commit and Push
```bash
git add .
git commit -m "feat: add restaurant dashboard feature"
git push origin feature/your-feature-name
```

### 6. Create a Pull Request
- Use the PR template
- Link related issues
- Request reviews from maintainers

## 🏗️ Project Structure

```
surplus-supper/
├── backend/                 # Go backend API
│   ├── api/                # HTTP handlers
│   ├── middleware/         # Authentication & CORS
│   ├── userService/        # User business logic
│   ├── restaurantService/  # Restaurant business logic
│   ├── db/                 # Database migrations
│   └── main.go            # Application entry point
├── frontend-next/          # Next.js frontend
│   ├── src/
│   │   ├── app/           # App router pages
│   │   ├── components/    # React components
│   │   ├── lib/           # Utility functions
│   │   └── types/         # TypeScript types
│   └── public/            # Static assets
└── docker-compose.yml     # Development environment
```

## 📝 Coding Standards

### Backend (Go)
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use meaningful variable and function names
- Add comments for exported functions
- Handle errors properly
- Write unit tests for business logic

### Frontend (Next.js/React)
- Use TypeScript for type safety
- Follow React best practices
- Use functional components with hooks
- Implement proper error boundaries
- Write component tests with Jest/React Testing Library

### Database
- Use meaningful table and column names
- Add proper indexes for performance
- Include foreign key constraints
- Write migration scripts for schema changes

## 🧪 Testing

### Backend Testing
```bash
cd backend
go test ./... -v
go test -cover ./...
```

### Frontend Testing
```bash
cd frontend-next
npm run test
npm run test:coverage
```

## 📚 Documentation

- Update README.md for new features
- Add inline code comments
- Update API documentation
- Include setup instructions for new dependencies

## 🏷️ Issue Labels

- `bug` - Something isn't working
- `enhancement` - New feature or request
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention is needed
- `documentation` - Improvements or additions to documentation
- `frontend` - Frontend-related changes
- `backend` - Backend-related changes
- `database` - Database schema or migration changes

## 🎯 Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): description

feat(auth): add JWT authentication system
fix(api): resolve 404 error in restaurant endpoint
docs(readme): update installation instructions
test(user): add unit tests for user service
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

## 🤝 Code Review Process

1. **Automated Checks**: CI/CD pipeline runs tests and linting
2. **Review Request**: Assign reviewers to your PR
3. **Feedback**: Address review comments
4. **Approval**: Get approval from maintainers
5. **Merge**: PR is merged to main branch

## 🐛 Reporting Bugs

Use the [Bug Report Template](.github/ISSUE_TEMPLATE/bug_report.md) and include:
- Clear description of the bug
- Steps to reproduce
- Expected vs actual behavior
- Environment details
- Screenshots if applicable

## 💡 Suggesting Features

Use the [Feature Request Template](.github/ISSUE_TEMPLATE/feature_request.md) and include:
- Clear description of the feature
- Problem it solves
- Proposed solution
- Technical requirements

## 📞 Getting Help

- **Issues**: Create an issue for bugs or feature requests
- **Discussions**: Use GitHub Discussions for questions
- **Documentation**: Check the README and code comments

## 🎉 Recognition

Contributors will be recognized in:
- Project README
- Release notes
- GitHub contributors page

## 📄 License

By contributing to Surplus Supper, you agree that your contributions will be licensed under the same license as the project.

---

Thank you for contributing to Surplus Supper! 🌟
